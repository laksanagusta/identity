#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_KEY="56c290ad131b1f3e3131059c6c33ff46be0cff5cab3673de2bf2c1d81798b1d8"
BASE_URL="https://identity.marvcore.com"

echo -e "${YELLOW}=== Testing External API Endpoint ===${NC}\n"

# Test 1: URL dengan spasi (AKAN ERROR)
echo -e "${YELLOW}Test 1: URL dengan spasi (expected to fail)${NC}"
echo "URL: ${BASE_URL}/api/external/v1/users?page=1&limit=1&employee_id=in 200206062025062005"
curl -i -X GET "${BASE_URL}/api/external/v1/users?page=1&limit=1&employee_id=in 200206062025062005" \
  -H "x-api-key: $API_KEY" 2>&1 | head -20
echo -e "\n"

# Test 2: URL dengan proper encoding (SEHARUSNYA BERHASIL)
echo -e "${GREEN}Test 2: URL dengan proper URL encoding${NC}"
ENCODED_URL="${BASE_URL}/api/external/v1/users?page=1&limit=1&employee_id=in%20200206062025062005"
echo "URL: $ENCODED_URL"
curl -i -X GET "$ENCODED_URL" \
  -H "x-api-key: $API_KEY" 2>&1 | head -20
echo -e "\n"

# Test 3: URL tanpa operator "in" (SEHARUSNYA BERHASIL)
echo -e "${GREEN}Test 3: URL tanpa operator 'in'${NC}"
echo "URL: ${BASE_URL}/api/external/v1/users?page=1&limit=1&employee_id=200206062025062005"
curl -i -X GET "${BASE_URL}/api/external/v1/users?page=1&limit=1&employee_id=200206062025062005" \
  -H "x-api-key: $API_KEY" 2>&1 | head -20
echo -e "\n"

# Test 4: Basic list tanpa filter
echo -e "${GREEN}Test 4: Basic list tanpa employee_id filter${NC}"
echo "URL: ${BASE_URL}/api/external/v1/users?page=1&limit=1"
curl -i -X GET "${BASE_URL}/api/external/v1/users?page=1&limit=1" \
  -H "x-api-key: $API_KEY"
echo -e "\n"

# Test 5: Tanpa API Key (expected to fail with 401)
echo -e "${YELLOW}Test 5: Tanpa x-api-key header (expected to fail with 401)${NC}"
echo "URL: ${BASE_URL}/api/external/v1/users?page=1&limit=1"
curl -i -X GET "${BASE_URL}/api/external/v1/users?page=1&limit=1"
echo -e "\n"
