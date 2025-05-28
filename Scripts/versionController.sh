#!/bin/bash
# Fetch remote version
remote_url="https://raw.githubusercontent.com/AnkanSaha/ContainDB/main/VERSION"
remote_version=$(curl -s "$remote_url")
if [[ -z "$remote_version" ]]; then
  echo "Error: Unable to fetch remote version."
  exit 1
fi
echo "Current GitHub version: $remote_version"

# Read local version
local_version=$(cat "$(dirname "$0")/../VERSION" 2>/dev/null || echo "0.0.0")
echo "Local version: $local_version"

# Compare versions: returns 0 if first > second
ver_gt() {
  local IFS=.
  local i ver1=($1) ver2=($2)
  for ((i=${#ver1[@]}; i<${#ver2[@]}; i++)); do ver1[i]=0; done
  for ((i=${#ver2[@]}; i<${#ver1[@]}; i++)); do ver2[i]=0; done
  for ((i=0; i<${#ver1[@]}; i++)); do
    if ((10#${ver1[i]} > 10#${ver2[i]})); then return 0; fi
    if ((10#${ver1[i]} < 10#${ver2[i]})); then return 1; fi
  done
  return 1
}

# Exit if local is already ahead of remote
if ver_gt "$local_version" "$remote_version"; then
  echo "Local version ($local_version) is ahead of remote ($remote_version). Skipping update."
  exit 0
fi

# Prompt for new version
read -p "Enter new version: " new_version

# Validate version format
if ! [[ "$new_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Invalid version format. Use X.Y.Z (e.g., 1.0.0)."
  exit 1
fi

# Update local version file
echo "$new_version" > "$(dirname "$0")/../VERSION"

echo "Local version updated to $new_version."