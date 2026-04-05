import sys
from time import sleep

from matrix import main_img, receive_img, send_img, shutdown_img, test_img
from sense_hat import SenseHat

from commands import decode_command, encode_response

args = sys.argv[1:]

if len(args) == 0:
    print(encode_response("main", "not_ok", "bad_args"))
    exit(0)

cmd = ""
val = ""

try:
    cmd, val = decode_command(args[0])
except Exception:
    print(encode_response("main", "not_ok", "bad_args"))
    exit(0)

sense = SenseHat()


def animate_down(image, delay=0.02):
    render = [[0, 0, 0] for _ in range(64)]
    for y in range(8):
        part_idx = (y + 1) * 8
        render[:part_idx] = image[:part_idx]
        sense.set_pixels(render)
        sleep(delay)


def animate_up(image, delay=0.02):
    render = [[0, 0, 0] for _ in range(64)]
    for y in range(0, 8):
        part_idx = (y + 1) * 8
        render[-part_idx:] = image[-part_idx:]
        sense.set_pixels(render)
        sleep(delay)


def animate_left(image, delay=0.02):
    render = [[0, 0, 0] for _ in range(64)]
    image_rows = [[[0, 0, 0] for _ in range(8)] for _ in range(8)]
    # we need to convert the image from 64 elems to 8 rows of 8 pixels
    for y in range(8):
        part_idx = (y + 1) * 8
        image_rows[y] = image[part_idx - 8 : part_idx]
    for x in range(8):
        col = []
        for j in range(8):
            col.append(image_rows[j][x])
        for p in range(8):
            render[p * 8 + x] = col[p]
        sense.set_pixels(render)
        sleep(delay)


def animate_right(image, delay=0.02):
    render = [[0, 0, 0] for _ in range(64)]
    image_rows = [[[0, 0, 0] for _ in range(8)] for _ in range(8)]
    for y in range(8):
        part_idx = (y + 1) * 8
        image_rows[y] = image[part_idx - 8 : part_idx]
    cols = []
    for x in range(8):
        col = []
        for j in range(8):
            col.append(image_rows[j][x])
        cols.append(col)
    cols.reverse()
    for x in range(8):
        col = cols[x]
        for p in range(8):
            render[p * 8 + (7 - x)] = col[p]
        sense.set_pixels(render)
        sleep(delay)


images = {
    "rx": receive_img,
    "tx": send_img,
    "test": test_img,
    "main": main_img,
    "shutdown": shutdown_img,
}

animate = {
    "down": animate_down,
    "up": animate_up,
    "left": animate_left,
    "right": animate_right,
}

sense.clear()

if cmd == "clear":
    print(encode_response("main", "ok", "cleared"))
    exit(0)

img = images.get(cmd)

if img is None:
    print(encode_response("main", "not_ok", "bad_cmd"))
    exit(0)

val_parts = val.split(",")

if len(val_parts) < 1:
    print(encode_response("main", "not_ok", "bad_val"))
    exit(0)

anim_type = val_parts[0]

delay = 0.02

if len(val_parts) >= 2:
    if len(val_parts) > 2:
        print(encode_response("main", "not_ok", "bad_val"))
        exit(0)
    try:
        delay = float(val_parts[1])
    except Exception:
        print(encode_response("main", "not_ok", "bad_val"))
        exit(0)

anim = animate.get(anim_type)

if anim is None:
    print(encode_response("main", "not_ok", "bad_val"))
    exit(0)

anim(img, delay)
print(encode_response("main", "ok", "done"))
