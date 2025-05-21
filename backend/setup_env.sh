#!/bin/bash

# Colors for terminal output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# File to store environment variables
ENV_FILE=".env"

echo -e "${BLUE}=== Reddit API Credentials Setup ===${NC}"
echo "This script will help you set up the Reddit API credentials required for the Reddit Summarizer app."
echo ""

# Check if .env file already exists
if [ -f "$ENV_FILE" ]; then
  echo -e "${YELLOW}Existing $ENV_FILE file found.${NC}"
  echo ""
  
  # Show current values (masked)
  if grep -q "redditID" "$ENV_FILE"; then
    ID_EXISTS=true
    echo -e "Reddit App ID: ${GREEN}[Already set]${NC}"
  else
    ID_EXISTS=false
    echo -e "Reddit App ID: ${RED}[Not set]${NC}"
  fi
  
  if grep -q "redditSecret" "$ENV_FILE"; then
    SECRET_EXISTS=true
    echo -e "Reddit App Secret: ${GREEN}[Already set]${NC}"
  else
    SECRET_EXISTS=false
    echo -e "Reddit App Secret: ${RED}[Not set]${NC}"
  fi
  
  if grep -q "redditUser" "$ENV_FILE"; then
    USER_EXISTS=true
    echo -e "Reddit Username: ${GREEN}[Already set]${NC}"
  else
    USER_EXISTS=false
    echo -e "Reddit Username: ${RED}[Not set]${NC}"
  fi
  
  if grep -q "redditPassword" "$ENV_FILE"; then
    PASS_EXISTS=true
    echo -e "Reddit Password: ${GREEN}[Already set]${NC}"
  else
    PASS_EXISTS=false
    echo -e "Reddit Password: ${RED}[Not set]${NC}"
  fi
  
  echo ""
  read -p "Do you want to update these credentials? (y/n): " UPDATE_CREDS
  
  if [[ $UPDATE_CREDS != "y" && $UPDATE_CREDS != "Y" ]]; then
    echo -e "${GREEN}Keeping existing credentials. Setup complete!${NC}"
    echo "Run 'source $ENV_FILE' to load these credentials into your current shell."
    exit 0
  fi
else
  echo -e "${YELLOW}No existing credentials found. Let's set them up!${NC}"
  echo ""
  # Create a new .env file
  touch "$ENV_FILE"
  ID_EXISTS=false
  SECRET_EXISTS=false
  USER_EXISTS=false
  PASS_EXISTS=false
fi

# Instructions for getting Reddit API credentials
echo -e "${BLUE}===== How to get Reddit API Credentials =====${NC}"
echo "1. Go to https://www.reddit.com/prefs/apps"
echo "2. Scroll down and click 'create another app...'"
echo "3. Fill in the form:"
echo "   - Name: RedditSummarizer (or any name you prefer)"
echo "   - Select 'script' as the application type"
echo "   - Description: App for summarizing Reddit comments"
echo "   - About URL: You can leave this blank"
echo "   - Redirect URI: http://localhost:8080"
echo "4. Click 'create app'"
echo "5. The 'App ID' is displayed under the app name"
echo "6. The 'App Secret' is displayed next to 'secret'"
echo ""
echo -e "${YELLOW}Note: You'll also need your Reddit username and password${NC}"
echo ""

# Function to update a specific credential in the .env file
update_credential() {
  local name=$1
  local var_name=$2
  local exists=$3
  local is_secret=$4
  
  if [ "$exists" = true ]; then
    if [ "$is_secret" = true ]; then
      read -p "Update $name? (leave blank to keep existing, or enter new value): " value
    else
      read -p "Update $name? (leave blank to keep existing, or enter new value): " value
    fi
    
    if [ -n "$value" ]; then
      # Remove existing line if it exists
      sed -i.bak "/^export $var_name=/d" "$ENV_FILE" && rm -f "$ENV_FILE.bak"
      # Add new value
      echo "export $var_name=\"$value\"" >> "$ENV_FILE"
      echo -e "${GREEN}$name updated!${NC}"
    else
      echo -e "${YELLOW}Keeping existing $name.${NC}"
    fi
  else
    if [ "$is_secret" = true ]; then
      read -s -p "Enter your $name: " value
      echo ""
    else
      read -p "Enter your $name: " value
    fi
    
    if [ -n "$value" ]; then
      echo "export $var_name=\"$value\"" >> "$ENV_FILE"
      echo -e "${GREEN}$name set!${NC}"
    else
      echo -e "${RED}$name is required for the app to function properly.${NC}"
    fi
  fi
}

# Update credentials
update_credential "Reddit App ID" "redditID" $ID_EXISTS false
update_credential "Reddit App Secret" "redditSecret" $SECRET_EXISTS false
update_credential "Reddit Username" "redditUser" $USER_EXISTS false
update_credential "Reddit Password" "redditPassword" $PASS_EXISTS true

# Make the file executable
chmod +x "$ENV_FILE"

echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo "To load these credentials into your current shell, run:"
echo -e "${YELLOW}    source $ENV_FILE${NC}"
echo ""
echo "To start the Reddit server with these credentials, run:"
echo -e "${YELLOW}    source $ENV_FILE && cd cmd/reddit-server && go run main.go${NC}"

