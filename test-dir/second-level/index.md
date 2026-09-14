#!/bin/env python3

import sys
import os

save_file = os.path.join(os.path.dirname(__file__), "save.dat")

# logins don't expire, don't use in production!
logins = {}

try: # might not exist
    with open(save_file, "r+") as f:
        for line in f.readlines():
            line = line.strip()
            if line != "":
                parts = line.split(" ", 1)
                logins[parts[0]] = parts[1]
except: pass

# on write
if len(sys.argv) == 3:
    logins[sys.argv[1]] = sys.argv[2]

    with open(save_file, "w") as f:
        for pid in logins:
            f.write(f"{pid} {logins[pid]}\n")

    print("you are now logged in as:", sys.argv[2])
    print("read again to see the changes!")
    exit(0)


# on read
print("# Dynamic file!")
print("This file has been dynamically generated")
print("\nyour PID is:", sys.argv[1], "\n")

if sys.argv[1] not in logins:
    print("You aren't logged in yet!\nWrite your username into this file to log in!")
else:
    print("You are logged int as:", logins[sys.argv[1]])
