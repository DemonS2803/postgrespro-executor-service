#!/bin/bash

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
echo "Factorial of $random_number is: $result"