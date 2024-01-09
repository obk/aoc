def calculate_floor(data):
    floor = 0
    for char in data:
        if char == '(':
            floor += 1
        elif char == ')':
            floor -= 1
    return floor

def find_basement_position(data):
    floor = 0
    position = 0
    for char in data:
        position += 1
        if char == '(':
            floor += 1
        elif char == ')':
            floor -= 1
        if floor == -1:
            return position
    return None

data = open("inputs/day1.txt").read()

print('Part 1:', calculate_floor(data))
print('Part 2:', find_basement_position(data))