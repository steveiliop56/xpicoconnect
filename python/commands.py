def encode_command(command, value):
    res = command
    res += ":"
    if not value:
        res += "null"
    else:
        res += str(value)
    return res


def decode_command(command):
    parts = command.split(":")
    if len(parts) != 2:
        raise ValueError("Invalid command format")
    return parts[0], parts[1].strip()


def decode_response(response):
    parts = response.split(":")
    if len(parts) != 3:
        raise ValueError("Invalid response format")
    if parts[1] != "ok":
        raise ValueError("Response indicates failure")
    return parts[2].strip()


def encode_response(command, status, value):
    res = command
    res += ":"
    res += status
    res += ":"
    if not value:
        res += "null"
    else:
        res += str(value)
    return res
