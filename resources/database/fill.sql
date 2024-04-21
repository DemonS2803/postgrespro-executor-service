

INSERT INTO commands (code, description, created_at)
VALUES
    ('echo $RANDOM', 'test', now()),
    ('#!/bin/bash

# Greet the user
echo "Hello! Nice to meet you."', 'second one', now()),
    ('#!/bin/bash
random_number=$((RANDOM % 100 + 1))
if [ $((random_number % 2)) -eq 0 ]; then
    echo "$random_number is even."
else
    echo "$random_number is odd."
fi', 'bigger', now()),
    ('#!/bin/bash
factorial() {
    if [ $1 -eq 0 ]; then
        echo 1
    else
        echo $(( $1 * $(factorial $(( $1 - 1 ))) ))
    fi
}
random_number=$((RANDOM % 10 + 1))
echo "Random number: $random_number"
result=$(factorial $random_number)
echo "Factorial of $random_number is: $result"', 'randomized factorial', now()),
    ('ping 8.8.8.8', 'ping ggoooooogle', now());


