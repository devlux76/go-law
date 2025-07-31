import os
import re
import asyncio
import bz2
import shutil
from pathlib import Path
from urllib.parse import urljoin

from playwright.async_api import async_playwright
import requests

BULK_URL = "https://com-courtlistener-storage.s3-us-west-2.amazonaws.com/list.html?prefix=bulk-data/"
BULK_DIR = "/srv/bulk-data"

ARCHIVE_REGEX = re.compile(r"^(?P<type>.+)-(?P<date>\d{4}-\d{2}-\d{2})\.csv\.bz2$")

def create_or_update_symlink(target, link_name):
    """Create or update a symlink safely."""
    link = Path(link_name)
    try:
        if link.is_symlink() or link.exists():
            link.unlink()
        link.symlink_to(Path(target).name)
    except Exception as e:
        print(f"Failed to create symlink {link_name} -> {target}: {e}")

def get_latest_file_for_type(files, file_type):
    """Return the latest dated file for a given type."""
    pattern = re.compile(rf"^{re.escape(file_type)}-(\d{{4}}-\d{{2}}-\d{{2}})\.csv$")
    dated_files = []
    for f in files:
        m = pattern.match(f)
        if m:
            dated_files.append((m.group(1), f))
    if not dated_files:
        return None
    # Sort by date descending
    dated_files.sort(reverse=True)
    return dated_files[0][1]

def cleanup_old_files(files, current_symlinks):
    """Remove files not referenced by any symlink."""
    referenced = set()
    for symlink in current_symlinks:
        try:
            target = os.readlink(symlink)
            referenced.add(target)
        except OSError:
            continue
    for f in files:
        if f not in referenced and not f.endswith('.sql'):
            try:
                os.remove(os.path.join(BULK_DIR, f))
            except Exception as e:
                print(f"Failed to remove old file {f}: {e}")

def get_existing_csvs():
    csvs = {}
    for f in os.listdir(BULK_DIR):
        m = re.match(r"(.+)-(\d{4}-\d{2}-\d{2})\.csv$", f)
        if m:
            csvs[m.group(0)] = True
    return csvs

async def fetch_links():
    async with async_playwright() as p:
        browser = await p.chromium.launch()
        page = await browser.new_page()
        await page.goto(BULK_URL)
        await page.wait_for_selector("a")
        links = await page.eval_on_selector_all("a", "els => els.map(e => e.href)")
        await browser.close()
        return links

def parse_archives(links):
    latest = {}
    for link in links:
        fname = link.split("/")[-1]
        m = ARCHIVE_REGEX.match(fname)
        if not m:
            continue
        typ = m.group("type")
        date = m.group("date")
        if typ not in latest or date > latest[typ][0]:
            latest[typ] = (date, link, fname)
    return latest

def download_and_decompress(url, fname):
    bz2_path = os.path.join(BULK_DIR, fname)
    csv_path = bz2_path[:-4]
    # Download
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        with open(bz2_path, "wb") as f:
            shutil.copyfileobj(r.raw, f)
    # Decompress
    with bz2.open(bz2_path, "rb") as f_in, open(csv_path, "wb") as f_out:
        shutil.copyfileobj(f_in, f_out)
    os.remove(bz2_path)

async def main():
    os.makedirs(BULK_DIR, exist_ok=True)
    existing = get_existing_csvs()
    links = await fetch_links()
    latest = parse_archives(links)
    for typ, (date, url, fname) in latest.items():
        csv_name = f"{typ}-{date}.csv"
        if csv_name in existing:
            continue
        print(f"Downloading {fname} ...")
        download_and_decompress(url, fname)
        print(f"Decompressed {csv_name}")

    # Symlink management
    files = os.listdir(BULK_DIR)
    types = set()
    for f in files:
        m = re.match(r"(.+)-(\d{4}-\d{2}-\d{2})\.csv$", f)
        if m:
            types.add(m.group(1))
    symlinks = []
    for typ in types:
        latest_file = get_latest_file_for_type(files, typ)
        if latest_file:
            link_name = os.path.join(BULK_DIR, f"{typ}-latest.csv")
            create_or_update_symlink(os.path.join(BULK_DIR, latest_file), link_name)
            symlinks.append(link_name)
    cleanup_old_files(files, symlinks)

if __name__ == "__main__":
    asyncio.run(main())