import sys
from time import sleep
import os

if __name__ == '__main__':
    print("hello world!")
    print(os.getenv("BEGIN_YEAR"))
    print(os.getenv("END_YEAR"))
    sleep(3)
    sys.exit(0)
