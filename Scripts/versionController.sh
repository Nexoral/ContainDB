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
  local raw1 raw2 i ver1 ver2
  # Strip suffix from versions (ignore -beta or -stable)
  raw1="${1%%-*}"
  raw2="${2%%-*}"
  ver1=($raw1) ver2=($raw2)
  # Add version type selection using arrow keys
  echo "Select version type:"
  options=("Stable" "Beta")
  selected=0

  # ANSI codes for cursor movement and clearing
  tput civis # Hide cursor

  # Display function
  function show_menu() {
    for i in ${!options[@]}; do
      if [[ $i -eq $selected ]]; then
        echo -e "\033[7m> ${options[$i]}\033[0m" # Highlighted
      else
        echo "  ${options[$i]}"
      fi
    done
  }

  # Clear previous output and show menu
  function refresh() {
    tput cup 4 0
    tput ed
    show_menu
  }

  # Initial display
  show_menu

  # Handle key input
  while true; do
    read -rsn3 key
    case "$key" in
    $'\x1b[A') # Up arrow
      ((selected--))
      if ((selected < 0)); then selected=$((${#options[@]} - 1)); fi
      refresh
      ;;
    $'\x1b[B') # Down arrow
      ((selected++))
      if ((selected >= ${#options[@]})); then selected=0; fi
      refresh
      ;;
    "") # Enter key
      version_type="${options[$selected]}"
      break
      ;;
    esac
  done

  tput cnorm # Show cursor again
  echo -e "\nSelected: $version_type"
  for ((i = ${#ver1[@]}; i < ${#ver2[@]}; i++)); do ver1[i]=0; done
  for ((i = ${#ver2[@]}; i < ${#ver1[@]}; i++)); do ver2[i]=0; done
  for ((i = 0; i < ${#ver1[@]}; i++)); do
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
suffix=$(echo "$version_type" | tr '[:upper:]' '[:lower:]')
echo "${new_version}-${suffix}" >"$(dirname "$0")/../VERSION"
echo "Local version updated to ${new_version}-${suffix}"

# Update static version in main.go and Banner.go
main_go="$(dirname "$0")/../src/Root/main.go"
sed -i "s|^\([[:space:]]*\)VERSION :=.*|\1VERSION := \"${new_version}-${suffix}\"|" "$main_go"

banner_go="$(dirname "$0")/../src/Root/Banner.go"
sed -i "s|^\([[:space:]]*\)const Version =.*|\1const Version = \"${new_version}-${suffix}\"|" "$banner_go"

echo "Updated main.go and Banner.go with new version."

# Update version in README.md installation URLs
readme_file="$(dirname "$0")/../README.md"
sed -i "s|^\(\s*wget .*/releases/download/v\)[^/]*\(/containDB_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2${new_version}-${suffix}\3|" "$readme_file"
sed -i "s|^\(\s*sudo dpkg -i containDB_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2|" "$readme_file"

echo "Updated README.md with new version in installation URLs."
