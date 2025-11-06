#!/usr/bin/env python3
"""
brute_steg.py
Brute-force offset/interval combos for ./steg (byte mode).
Detects failure by the exact stderr message:
    Error: Sentinel not found (byte mode).
and treats a run as success when stderr does NOT contain that message and stdout is non-empty.

Adjust OFFSETS and INTERVALS as needed.
"""
import subprocess
import sys
import os
import time
from pathlib import Path

STEG_BIN = "./steg"       # path to your steg binary; change if needed
INPUT_FILE = "const.bin"  # input file to probe
OUT_DIR = Path("results") # where successful stdout blobs are saved
TIMEOUT = 5               # seconds per attempt (increase if your machine is slow)

# grid to try
OFFSETS = range(0, 256)                   # 0..255 by default
INTERVALS = [1,2,4,8,16,32,64,128,256]    # powers-of-two intervals

# failure string (exact text printed when wrong params are used)
FAILURE_PHRASE = "Error: Sentinel not found (byte mode)."

# Prepare output dir
OUT_DIR.mkdir(parents=True, exist_ok=True)

def run_steg(offset, interval):
    cmd = [STEG_BIN, "-r", "-B", f"-o{offset}", f"-i{interval}", f"-w{INPUT_FILE}"]
    try:
        proc = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=TIMEOUT)
    except subprocess.TimeoutExpired:
        return False, b"", b"__TIMEOUT__"
    return True, proc.stdout, proc.stderr

def is_failure(stderr_bytes):
    try:
        s = stderr_bytes.decode("utf-8", errors="ignore")
    except Exception:
        s = str(stderr_bytes)
    return FAILURE_PHRASE.lower() in s.lower()

def save_success(offset, interval, out_bytes, err_bytes):
    base = OUT_DIR / f"out_o{offset}_i{interval}"
    out_path = base.with_suffix(".bin")
    err_path = base.with_suffix(".err.txt")
    with open(out_path, "wb") as f:
        f.write(out_bytes)
    with open(err_path, "wb") as f:
        f.write(err_bytes)
    return out_path, err_path

def printable_preview(b, length=512):
    try:
        txt = b.decode("utf-8", errors="ignore")
    except Exception:
        txt = str(b)
    # show first 'length' bytes/chars, but keep it safe for terminal
    return txt[:length].replace("\r", "")

def main():
    if not Path(STEG_BIN).exists():
        print(f"ERROR: steg binary not found at {STEG_BIN}. Edit STEG_BIN variable in the script.")
        sys.exit(1)
    if not Path(INPUT_FILE).exists():
        print(f"ERROR: input file not found: {INPUT_FILE}")
        sys.exit(1)

    total = 0
    found = []
    start = time.time()

    # iterate grid
    for off in OFFSETS:
        for iv in INTERVALS:
            total += 1
            ok, sout, serr = run_steg(off, iv)
            if not ok:
                # timeout happened
                print(f"[TIMEOUT] o={off} i={iv}")
                continue

            # if stderr contains the failure phrase -> skip
            if is_failure(serr):
                # explicit failure
                # (we intentionally do not save these)
                continue

            # If stderr is empty or doesn't have the failure phrase, check stdout content:
            if sout and len(sout) > 0:
                # treat this as a success; save and show a short preview
                out_path, err_path = save_success(off, iv, sout, serr)
                preview = printable_preview(sout, 768)
                print(f"[SUCCESS] o={off} i={iv} -> saved {out_path} (stderr -> {err_path})")
                # print small preview so you can quickly inspect
                print("----- preview -----")
                print(preview)
                print("-------------------\n")
                found.append((off, iv, str(out_path)))
            else:
                # stdout empty -> probably failure / nothing extracted
                # there are cases where steg prints nothing but creates a file when redirecting:
                # since we're capturing stdout here, empty stdout means nothing to save.
                continue

    elapsed = time.time() - start
    print("Done.")
    print(f"Total combos tried: {total}")
    print(f"Successful extractions saved: {len(found)}")
    for o, i, path in found:
        print(f"  offset={o} interval={i} -> {path}")
    print(f"Elapsed {elapsed:.1f}s")

if __name__ == "__main__":
    main()
