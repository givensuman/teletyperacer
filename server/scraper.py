import requests
from bs4 import BeautifulSoup
import sqlite3
import time
import re


def init_db():
    conn = sqlite3.connect("texts.db")
    c = conn.cursor()
    c.execute(
        """CREATE TABLE IF NOT EXISTS texts (id INTEGER PRIMARY KEY, text TEXT, source TEXT)"""
    )
    conn.commit()
    return conn


def scrape_text(id):
    url = f"https://typeracerdata.com/text?id={id}"
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
    except requests.RequestException as e:
        print(f"Error fetching {id}: {e}")
        return None, None

    soup = BeautifulSoup(response.content, "html.parser")

    # Find the text: it's in a <p> tag after the h1 with "Text #id"
    h1 = soup.find("h1", string=re.compile(r"Text #\d+"))
    if h1:
        p_text = h1.find_next("p")
        if p_text:
            text = p_text.get_text(strip=True)
            p_source = p_text.find_next("p")
            source = p_source.get_text(strip=True) if p_source else ""
            return text, source

    # Fallback: find any p with long text
    paragraphs = soup.find_all("p")
    for i, p in enumerate(paragraphs):
        t = p.get_text(strip=True)
        if len(t) > 100 and not t.startswith("View"):
            text = t
            source = (
                paragraphs[i + 1].get_text(strip=True)
                if i + 1 < len(paragraphs)
                else ""
            )
            return text, source

    return None, None


def store_text(conn, id, text, source):
    c = conn.cursor()
    c.execute(
        "INSERT OR REPLACE INTO texts (id, text, source) VALUES (?, ?, ?)",
        (id, text, source),
    )
    conn.commit()


def scrape_and_store(start, end):
    conn = init_db()
    for id in range(start, end + 1):
        text, source = scrape_text(id)
        if text:
            store_text(conn, id, text, source)
            print(f"Stored text {id}: {source}")
        else:
            print(f"No text for {id}")
        time.sleep(1)  # Be nice
    conn.close()


if __name__ == "__main__":
    scrape_and_store(3640650, 3640660)  # Example range around the known ID
