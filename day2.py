def part1(data):
    sum = 0
    for x in data:
        l = int(x[0])
        w = int(x[1])
        h = int(x[2])
        sum += 2*l*w + 2*w*h + 2*h*l + min(l*w, w*h, h*l)
    return sum

def part2(data):
    sum = 0
    for x in data:
        l = int(x[0])
        w = int(x[1])
        h = int(x[2])
        sum += l*w*h + min(2*l+2*w, 2*w+2*h, 2*h+2*l)
    return sum

data = open("inputs/day2.txt").read().splitlines()
data = [x.split('x') for x in data]

print('Part 1:', part1(data))
print('Part 2:', part2(data))
