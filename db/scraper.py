import requests
from bs4 import BeautifulSoup
import sqlite3
import time
import re


def init_db():
    conn = sqlite3.connect("texts.db")
    c = conn.cursor()
    c.execute(
        """CREATE TABLE IF NOT EXISTS texts (id INTEGER PRIMARY KEY, text TEXT, attribution TEXT)"""
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
        return None

    soup = BeautifulSoup(response.content, "html.parser")

    res = {
        "text": None,
        "attribution": None,
    }

    # Find the text: it's in a <p> tag after the h1 with "Text #id"
    h1 = soup.find("h1", string=re.compile(r"Text #\d+"))
    if h1:
        p = h1.find_next("p")
        if p:
            res["text"] = p.get_text(strip=True)
        p = p.find_next("p")
        if p:
            res["attribution"] = p.get_text(strip=True)

        if res["text"] is not None and res["attribution"] is not None:
            print(res["text"], "\n\t- ", res["attribution"])
            return res

    return None


def store_text(conn, id, res):
    c = conn.cursor()
    c.execute(
        "INSERT OR REPLACE INTO texts (id, text, attribution) VALUES (?, ?, ?)",
        (id, res["text"], res["attribution"]),
    )
    conn.commit()


def scrape_and_store(start, end):
    conn = init_db()
    for id in range(start, end + 1):
        text = scrape_text(id)
        if text:
            store_text(conn, id, text)
            print(f"Stored text {id}")
        else:
            print(f"No text for {id}")
        time.sleep(1)  # Be nice
    conn.close()


if __name__ == "__main__":
    scrape_and_store(3640650, 3640660)  # Example range around the known ID
