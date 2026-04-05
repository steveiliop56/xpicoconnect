K = [0, 0, 0]        # Off
G = [0, 200, 0]      # Green 
C = [0, 200, 255]    # Cyan   
W = [200, 200, 200]  # White  
Y = [255, 200, 0]    # Yellow 
R = [255, 50, 50]    # Red   

# Plane landing
receive_img = [
    K, K, K, K, K, K, K, K,
    K, K, K, K, G, K, K, K,
    K, K, K, K, G, K, K, K,
    K, K, G, K, G, K, G, K,
    K, K, K, G, G, G, K, K,
    K, K, K, K, G, K, K, K,
    K, K, K, K, G, K, K, K,
    K, K, K, K, K, K, K, K,
]

# Plane taking off
send_img = [
    K, K, K, K, K, K, K, K,
    K, K, K, K, C, K, K, K,
    K, K, K, K, C, K, K, K,
    K, K, K, C, C, C, K, K,
    K, K, C, K, C, K, C, K,
    K, K, K, K, C, K, K, K,
    K, K, K, K, C, K, K, K,
    K, K, K, K, K, K, K, K,
]

# Propeller
main_img = [
    K, K, K, K, K, K, K, K,
    K, K, K, K, W, K, K, K,
    K, K, K, K, W, K, K, K,
    K, K, W, W, W, W, W, K,
    K, K, K, K, W, K, K, K,
    K, K, K, K, W, K, K, K,
    K, K, K, K, K, K, K, K,
    K, K, K, K, K, K, K, K,
]

# Paper airplane
test_img = [
    K, K, K, K, K, K, K, K,
    K, K, K, Y, K, K, K, K,
    K, K, K, K, Y, K, K, K,
    K, K, Y, Y, Y, Y, Y, K,
    K, K, K, K, Y, K, K, K,
    K, K, K, Y, K, K, K, K,
    K, K, K, K, K, K, K, K,
    K, K, K, K, K, K, K, K,
]

# Parachute
shutdown_img = [
    K, K, K, K, K, K, K, K,
    K, K, R, R, R, R, K, K,
    K, R, K, K, K, K, R, K,
    K, K, R, K, K, R, K, K,
    K, K, K, R, R, K, K, K,
    K, K, K, R, R, K, K, K,
    K, K, K, K, K, K, K, K,
    K, K, K, K, K, K, K, K,
]
