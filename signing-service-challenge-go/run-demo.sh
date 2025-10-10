#!/bin/bash

# Signature Service Demo Script
# This script demonstrates the signing service API functionality

set -e

BASE_URL="http://localhost:8080/api/v0"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "\n${BOLD}${BLUE}=== $1 ===${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_step() {
    echo -e "\n${YELLOW}$1${NC}"
}

# Function to make HTTP requests and pretty print JSON
make_request() {
    local method=$1
    local path=$2
    local data=$3
    
    if [ -z "$data" ]; then
        curl -s -X "$method" "$BASE_URL$path"
    else
        curl -s -X "$method" "$BASE_URL$path" \
            -H "Content-Type: application/json" \
            -d "$data"
    fi
}

# Function to extract device ID from JSON response
extract_device_id() {
    echo "$1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4
}

# Wait for server to be ready
wait_for_server() {
    print_step "Checking if server is running..."
    for i in {1..5}; do
        if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
            print_success "Server is ready"
            return 0
        fi
        sleep 1
    done
    echo "Error: Server not running. Please start the server first with: go run main.go"
    exit 1
}

# Main demo flow
main() {
    print_header "Signature Service Demo"
    
    wait_for_server
    
    # 1. Create RSA Device
    print_step "1. Creating RSA Device..."
    RSA_RESPONSE=$(make_request "POST" "/devices" '{"algorithm":"RSA","label":"Demo RSA Device"}')
    DEVICE_ID_RSA=$(extract_device_id "$RSA_RESPONSE")
    print_success "Created device: $DEVICE_ID_RSA (RSA)"
    echo "$RSA_RESPONSE" | jq '.' 2>/dev/null || echo "$RSA_RESPONSE"
    
    # 2. Create ECC Device
    print_step "2. Creating ECC Device..."
    ECC_RESPONSE=$(make_request "POST" "/devices" '{"algorithm":"ECC","label":"Demo ECC Device"}')
    DEVICE_ID_ECC=$(extract_device_id "$ECC_RESPONSE")
    print_success "Created device: $DEVICE_ID_ECC (ECC)"
    echo "$ECC_RESPONSE" | jq '.' 2>/dev/null || echo "$ECC_RESPONSE"
    
    # 3. List all devices
    print_step "3. Listing all devices..."
    LIST_RESPONSE=$(make_request "GET" "/devices")
    print_success "Devices retrieved:"
    echo "$LIST_RESPONSE" | jq '.' 2>/dev/null || echo "$LIST_RESPONSE"
    
    # 4. Sign first transaction with RSA
    print_step "4. Signing first transaction (RSA)..."
    SIGN1_RESPONSE=$(make_request "POST" "/signatures" "{\"device_id\":\"$DEVICE_ID_RSA\",\"data\":\"transaction_data_1\"}")
    print_success "Transaction signed:"
    echo "$SIGN1_RESPONSE" | jq '.' 2>/dev/null || echo "$SIGN1_RESPONSE"
    
    # 5. Sign second transaction with RSA (shows chaining)
    print_step "5. Signing second transaction (RSA) - demonstrating signature chaining..."
    SIGN2_RESPONSE=$(make_request "POST" "/signatures" "{\"device_id\":\"$DEVICE_ID_RSA\",\"data\":\"transaction_data_2\"}")
    print_success "Transaction signed with chaining:"
    echo "$SIGN2_RESPONSE" | jq '.' 2>/dev/null || echo "$SIGN2_RESPONSE"
    
    # 6. Sign third transaction
    print_step "6. Signing third transaction (RSA)..."
    SIGN3_RESPONSE=$(make_request "POST" "/signatures" "{\"device_id\":\"$DEVICE_ID_RSA\",\"data\":\"transaction_data_3\"}")
    print_success "Transaction signed:"
    echo "$SIGN3_RESPONSE" | jq '.' 2>/dev/null || echo "$SIGN3_RESPONSE"
    
    # 7. Get device to check counter
    print_step "7. Checking device state (counter should be 3)..."
    DEVICE_RESPONSE=$(make_request "GET" "/devices/$DEVICE_ID_RSA")
    print_success "Device state:"
    echo "$DEVICE_RESPONSE" | jq '.' 2>/dev/null || echo "$DEVICE_RESPONSE"
    
    # 8. Sign with ECC device
    print_step "8. Signing transaction with ECC device..."
    SIGN_ECC_RESPONSE=$(make_request "POST" "/signatures" "{\"device_id\":\"$DEVICE_ID_ECC\",\"data\":\"ecc_transaction_data\"}")
    print_success "Transaction signed (ECC):"
    echo "$SIGN_ECC_RESPONSE" | jq '.' 2>/dev/null || echo "$SIGN_ECC_RESPONSE"
    
    # Summary
    print_header "Demo Complete"
    echo "Key Observations:"
    echo "- Each signature contains the previous signature (chaining)"
    echo "- Signature counter increments monotonically"
    echo "- Different algorithms (RSA, ECC) produce different signature formats"
    echo "- All operations are thread-safe and concurrent-ready"
    echo ""
    echo "Note: Install 'jq' for better JSON formatting (optional)"
}

# Run the demo
main

