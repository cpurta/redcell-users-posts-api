#!/bin/bash

# API Testing Script
# Tests user endpoints: GET /users, POST /users, GET /users/:id, PUT /users/:id, DELETE /users/:id

# Configuration
# Set API_URL to either local address or the provided IP
# API_URL="http://localhost:8080"  # Uncomment for local testing
API_URL="http://34.60.24.109"      # Remote server

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
  local status=$1
  local message=$2
  
  if [ "$status" = "PASS" ]; then
    echo -e "${GREEN}[PASS]${NC} $message"
  elif [ "$status" = "FAIL" ]; then
    echo -e "${RED}[FAIL]${NC} $message"
  else
    echo -e "${YELLOW}[$status]${NC} $message"
  fi
}

# Function to validate HTTP status code
check_status() {
  local expected=$1
  local received=$2
  local message=$3
  
  if [ "$received" -eq "$expected" ]; then
    print_status "PASS" "$message (Status: $received)"
    return 0
  else
    print_status "FAIL" "$message (Expected: $expected, Got: $received)"
    return 1
  fi
}

echo "===== API TESTING SCRIPT ====="
echo "Target API: $API_URL"
echo "============================="

# Step 1: Get all users
echo -e "\n--- Step 1: Get all users ---"
users_response=$(curl -s -w "\n%{http_code}" $API_URL/users)
users_status=$(echo "$users_response" | tail -n1)
users_body=$(echo "$users_response" | sed '$d')

check_status 200 $users_status "GET /users"
echo "Current users in system:"
echo "$users_body" | jq '.' || echo "$users_body"

# Step 2: Create a new user
echo -e "\n--- Step 2: Create a new user ---"
timestamp=$(date +%s)
new_user=$(cat <<EOF
{
  "first_name": "Test",
  "last_name": "User${timestamp}",
  "email": "testuser${timestamp}@example.com"
}
EOF
)

echo "Creating user with data:"
echo "$new_user" | jq '.'

create_response=$(curl -s -w "\n%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$new_user" \
  $API_URL/users)

create_status=$(echo "$create_response" | tail -n1)
create_body=$(echo "$create_response" | sed '$d')

check_status 201 $create_status "POST /users"
echo "Created user:"
echo "$create_body" | jq '.' || echo "$create_body"

# Extract the user ID from the response
user_id=$(echo "$create_body" | jq -r '.id')

if [ -z "$user_id" ] || [ "$user_id" = "null" ]; then
  print_status "FAIL" "Failed to extract user ID from response"
  exit 1
fi

print_status "INFO" "New user ID: $user_id"

# Step 3: Verify the user was created
echo -e "\n--- Step 3: Verify user was created ---"
verify_response=$(curl -s -w "\n%{http_code}" $API_URL/users/$user_id)
verify_status=$(echo "$verify_response" | tail -n1)
verify_body=$(echo "$verify_response" | sed '$d')

check_status 200 $verify_status "GET /users/$user_id"
echo "Retrieved user details:"
echo "$verify_body" | jq '.' || echo "$verify_body"

# Step 4: Update the user
echo -e "\n--- Step 4: Update user ---"
updated_user=$(cat <<EOF
{
  "first_name": "Updated",
  "last_name": "User${timestamp}",
  "email": "updated${timestamp}@example.com"
}
EOF
)

echo "Updating user with data:"
echo "$updated_user" | jq '.'

update_response=$(curl -s -w "\n%{http_code}" \
  -X PUT \
  -H "Content-Type: application/json" \
  -d "$updated_user" \
  $API_URL/users/$user_id)

update_status=$(echo "$update_response" | tail -n1)
update_body=$(echo "$update_response" | sed '$d')

check_status 200 $update_status "PUT /users/$user_id"
echo "Updated user:"
echo "$update_body" | jq '.' || echo "$update_body"

# Step 5: Verify the update
echo -e "\n--- Step 5: Verify update ---"
verify_update_response=$(curl -s -w "\n%{http_code}" $API_URL/users/$user_id)
verify_update_status=$(echo "$verify_update_response" | tail -n1)
verify_update_body=$(echo "$verify_update_response" | sed '$d')

check_status 200 $verify_update_status "GET /users/$user_id (after update)"
echo "Retrieved updated user:"
echo "$verify_update_body" | jq '.' || echo "$verify_update_body"

# Check if first_name is updated correctly
first_name=$(echo "$verify_update_body" | jq -r '.first_name')
if [ "$first_name" = "Updated" ]; then
  print_status "PASS" "First name updated correctly"
else
  print_status "FAIL" "First name not updated (Expected: 'Updated', Got: '$first_name')"
fi

# Step 6: Delete the user
echo -e "\n--- Step 6: Delete user ---"
delete_response=$(curl -s -w "\n%{http_code}" -X DELETE $API_URL/users/$user_id)
delete_status=$(echo "$delete_response" | tail -n1)
delete_body=$(echo "$delete_response" | sed '$d')

check_status 204 $delete_status "DELETE /users/$user_id"
echo "Delete response:"
echo "$delete_body" | jq '.' || echo "$delete_body"

# Step 7: Verify deletion
echo -e "\n--- Step 7: Verify deletion ---"
verify_delete_response=$(curl -s -w "\n%{http_code}" $API_URL/users/$user_id)
verify_delete_status=$(echo "$verify_delete_response" | tail -n1)

# Here we expect 404 because the user should not exist anymore
check_status 404 $verify_delete_status "GET /users/$user_id (after deletion)"

echo -e "\n===== TEST SUMMARY ====="
print_status "INFO" "Tests completed against $API_URL"
echo "User flow tested: Create → Retrieve → Update → Delete"
echo "=========================="