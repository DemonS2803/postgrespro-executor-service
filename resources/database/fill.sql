

INSERT INTO commands (code, description, created_at)
VALUES
    ('echo $RANDOM', 'test', now()),
    ('#!/bin/bash

# Greet the user
echo "Hello! Nice to meet you."', 'second one', now()),
    ('#!/bin/bash

# Generate a random number between 1 and 100
random_number=$((RANDOM % 100 + 1))

# Check if the random number is even or odd
if [ $((random_number % 2)) -eq 0 ]; then
    echo "$random_number is even."
else
    echo "$random_number is odd."
fi', 'bigger', now()),
    ('#!/bin/bash

# Function to calculate factorial
factorial() {
    if [ $1 -eq 0 ]; then
        echo 1
    else
        echo $(( $1 * $(factorial $(( $1 - 1 ))) ))
    fi
}

# Generate a random number between 1 and 10
random_number=$((RANDOM % 10 + 1))

# Display the random number
echo "Random number: $random_number"

# Calculate the factorial of the random number
result=$(factorial $random_number)

# Display the result
echo "Factorial of $random_number is: $result"', 'randomized factorial', now());


