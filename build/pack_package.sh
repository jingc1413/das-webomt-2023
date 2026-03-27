#!/usr/bin/env sh
set -e

export NO_COLOR=0
export FORCE_COLOR=1

SYSLOG_ERROR="user.error"
SYSLOG_INFO="user.info"
SYSLOG_DEBUG="user.debug"
#error
SYSLOG_LEVEL_ERROR=3
#info
SYSLOG_LEVEL_INFO=6
#debug
SYSLOG_LEVEL_DEBUG=7
#debug2
SYSLOG_LEVEL_DEBUG_2=8
#debug3
SYSLOG_LEVEL_DEBUG_3=9

SYSLOG_LEVEL_DEFAULT=$SYSLOG_LEVEL_ERROR
#none
SYSLOG_LEVEL_NONE=0

__INTERACTIVE=""
if [ -t 1 ]; then
    __INTERACTIVE="1"
fi

__green() {
    if [ "${__INTERACTIVE}${NO_COLOR:-0}" = "10" -o "${FORCE_COLOR}" = "1" ]; then
        printf '\33[1;32m%b\33[0m' "$1"
        return
    fi
    printf -- "%b" "$1"
}

__red() {
    if [ "${__INTERACTIVE}${NO_COLOR:-0}" = "10" -o "${FORCE_COLOR}" = "1" ]; then
        printf '\33[1;31m%b\33[0m' "$1"
        return
    fi
    printf -- "%b" "$1"
}

_printargs() {
    _exitstatus="$?"
    if [ -z "$NO_TIMESTAMP" ] || [ "$NO_TIMESTAMP" = "0" ]; then
        printf -- "%s" "[$(date)] "
    fi
    if [ -z "$2" ]; then
        printf -- "%s" "$1"
    else
        printf -- "%s" "$1='$2'"
    fi
    printf "\n"
    # return the saved exit status
    return "$_exitstatus"
}

_syslog() {
    _exitstatus="$?"
    if [ "${SYS_LOG:-$SYSLOG_LEVEL_NONE}" = "$SYSLOG_LEVEL_NONE" ]; then
        return
    fi
    _logclass="$1"
    shift
    if [ -z "$__logger_i" ]; then
        if _contains "$(logger --help 2>&1)" "-i"; then
            __logger_i="logger -i"
        else
            __logger_i="logger"
        fi
    fi
    $__logger_i -t "$PROJECT_NAME" -p "$_logclass" "$(_printargs "$@")" >/dev/null 2>&1
    return "$_exitstatus"
}

_log() {
    [ -z "$LOG_FILE" ] && return
    _printargs "$@" >>"$LOG_FILE"
}

_info() {
    #_log "$@"
    if [ "${SYS_LOG:-$SYSLOG_LEVEL_NONE}" -ge "$SYSLOG_LEVEL_INFO" ]; then
        _syslog "$SYSLOG_INFO" "$@"
    fi
    _printargs "$@"
}

_err() {
    _syslog "$SYSLOG_ERROR" "$@"
    #_log "$@"
    if [ -z "$NO_TIMESTAMP" ] || [ "$NO_TIMESTAMP" = "0" ]; then
        printf -- "%s" "[$(date)] " >&2
    fi
    if [ -z "$2" ]; then
        __red "$1" >&2
    else
        __red "$1='$2'" >&2
    fi
    printf "\n" >&2
    return 1
}

PROGRESS_BAR_WIDTH=50  # progress bar length in characters
draw_progress_bar() {
  # Arguments: current value, max value, unit of measurement (optional)
  local __value=$1
  local __max=$2
  local __unit=${3:-""}  # if unit is not supplied, do not display it

  # Calculate percentage
  if (( $__max < 1 )); then __max=1; fi  # anti zero division protection
  local __percentage=$(( 100 - ($__max*100 - $__value*100) / $__max ))

  # Rescale the bar according to the progress bar width
  local __num_bar=$(( $__percentage * $PROGRESS_BAR_WIDTH / 100 ))

  # Draw progress bar
  printf "["
  for b in $(seq 1 $__num_bar); do printf "#"; done
  for s in $(seq 1 $(( $PROGRESS_BAR_WIDTH - $__num_bar ))); do printf " "; done
  printf "] $__percentage%% ($__value / $__max $__unit)\r"
}

dst_size=0
function calc_dst_size(){
    dir=$(pwd)
    if [ -n "$1" ]; then
            dir=$1
    fi
    dst_size=$(du -sb $dir | awk '{print $1}')
}
function show_bar(){
    file_size=$(du -sb $1 | awk '{print $1}')
    while true; do
        calc_dst_size $2
        uploaded_bytes=$dst_size
        draw_progress_bar "${uploaded_bytes}" ${file_size} "bytes"
        if [ $uploaded_bytes == $file_size ]; then break; fi
        sleep 0.1  # Wait before redrawing
    done
}

crypt_passwd="SunWave321"

prikey='''-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAqL8q15Nf/3noMGSH7fuQzzk7I0WKDaeIU6aOUNF0Hrj4+jGM
R62i6PwIOEMKwfdaf8wl6DmEuBIml/X95zgTXH0HWl2HSX3AYn2WYpnOkfmVlt38
3xEBTIFBqstTTKgTo3y39ngQKWLbbXTOIf5AUzm6SqQVzlLj79kdhA5FKaUyXSQj
FGATEULsx8A1RR8QY9j+Mr9/jbvJiNkeEIH2mor84mQoQ4/9SvtX09c3G7B87My+
YN7mwqgFvb+51YvRSvQZrexQOd5XsfJ8K8RkwbMpjf1ZD1IwBhqdZflmJW93tucm
DaSmMdWFSYw/ld+5dVkKp4VTrZ0Mycdq2rpTpwIDAQABAoIBAQCFrCG48YKrMxuu
MgTHjW8x6EkjcLOii2LfuOGlvbX5nFeGgfd41GBnpTaxstHzwwjKkaI5qk6rLZ2q
5O+D2RTlQBmyCMLhgZ8Dpi6Z3vMXCZnpD+J/pc3cl8u4CybqY35jFKn5sTRERI5D
z7J2lRzJkMI03HR1o43ycpACCvfYuz10g3eLoaY/+XtLcPGutYskNxVBWIUsU35l
oRnc4EyJp9GXEbmp6qgq7QaG95BJUymH4ZWpajywVQR1tgKBABYBPRKzuLAkdK3m
fKbWXwiwev1Yx9iCIvaKeVXwEu34d2fDx7YJWypsE3wSYRUCOyPLEyafrxeWGM/z
JH/1EIxhAoGBANg8P8WucH92wWL47U5sL6o/EVdC3BxU6t8QTZzSWr3gG7HhjMsx
o1qkaOgKNT+kqoHV73bnIOA9oqiKhvHZVg/VgiwF/JOqBis+vnQF3sUkVfvBDhm/
A7rLkRfcfbQqmkWXxey+RrMUXP3DRIdPgHZXIg9dIm9gYn16PMVRHWDXAoGBAMfH
SeqYjE71gu/7QJCQZuVfhf8iHi8Usnetpp0janDZieFmCpv3679RQyvceqJ47zbn
XzyVKjMGdvp+rkEuwZP8te8Xi8DlSYETzsOvG0Ocp817VAYHXN5IzbecWG61ELGZ
tyu5cxXq+zp8bU7XN/dAXPj6iu+tR8juZ+Q4KbmxAoGBAJwAM335kH3U/kXaEtJe
KLEpWOhcexIRjXUqfOPjusV7ll9purqpcBGnxWuKWco1kTGkV2Ov8b71mJEBh8OZ
qYY44fXCx+r8YgD2/k4UIhiWU1YSfIrEjPtJe778OtAtYClPCuL2j2kJCuGk6563
E5XM3Oy6o2o43jVqZ8du8sP3AoGAUafYQ4YUm1VoLnSwwCX1mL1BhLXKRo4cICSf
HE1UfUm9PJ67qWJdPnaVkQDuMyhhBjztxVJmzDJRQTO0P7o/ryA0sMQcADz3nUWe
VodlCm6me6tz/X7W65gwVWMXFjD14NGmd722F3hTNWWUPAxluK7bEX0+epaF58/i
MPrvuxECgYEAyz0Xhum8jIJay8Eb2qLlwYJ+PniSqeKswlQ8kI1hpXBidml934c0
kby4mUax1NWvaFvUobibrmqbNkrfG+u9JLndu8H3q10FtGOOXOniQzTkA10thmL2
7qT4BB4u7b7yfYvKBu5wCHXuNWLRzPNzbsqUzuWL8yBgG/gfXaO9fgw=
-----END RSA PRIVATE KEY-----
'''

function get_current_time() {
    local formatted_time=$(date +"%Y%m%d")
    echo "$formatted_time"
}

function file_digest() {
    local fileselectdir=$1
    local savepath=$2
    local filepathsort=()
    local filedict=()
    local size=0
   

    for dirpath in $(find "$fileselectdir" -type d); do
       if [[ "$dirpath" == "." || "$dirpath" == ".." ]]; then
            echo "break" $dirpath
            continue
       fi
        echo "dirpath:" $dirpath
        for filename in $(find "$dirpath" -type f); do
            if [[ "$filename" == "." || "$filename" == ".." ]]; then
                echo "break" $filename
                continue
            fi
            if [[ "$filename" != *"signaturefile"* ]]; then
                #file_path="$dirpath/$filename"
                file_path="$filename"
                echo "------" $file_path $filename
                size=$(stat -c %s "$file_path")
               filedict["$size"]="$file_path"
            fi
        done
    done

    # Sort the filedict
    local sorted_filedict=($(for key in "${!filedict[@]}"; do echo "$key ${filedict[$key]}"; done | sort -n -r -k1))

    for item in "${sorted_filedict[@]}"; do
        echo $item
        filepathsort+=("${item#* }")
    done

    #echo $sorted_filedict
    #echo $filepathsort
    # Save the sorted files to a bin file
    cat "${filepathsort[@]}" > "$savepath/save.bin"
    # Calculate hash
    #local hash_value=$(sha256sum "$savepath/save.bin" | awk '{print $1}')
    local hash_value=$(openssl dgst -sha256 -binary "$savepath/save.bin")
    # Remove the bin file
    rm "$savepath/save.bin"
    echo "$hash_value"
}

signaturefile() {
    local hash_value=$1
    #tmppath 签名文件路径
    local tmppath=$2
    local fileselectdir=$3

    if [[ $(ls "$tmppath") =~ "signaturefile" ]]; then
        rm "$tmppath/signaturefile"
    fi

    cd "$tmppath"

    local signature="simulated_signature"
    # Write the signature to the file
    mkdir -p "$fileselectdir"
    echo "$signature" > "$fileselectdir/signaturefile"
}

create_signature_bak() {
    hash_value=$1
    tmppath=$2
    fileselectdir=$3
    prikey_file=$4
    
    if [ -f "${tmppath}/signaturefile" ]; then
        rm "${tmppath}/signaturefile"
    fi
    cd "${tmppath}"
    # use OpenSSL generate signature
    openssl rsautl -sign -inkey "${prikey_file}" -outform der -out "${tmppath}/signaturefile" <<< "${hash_value}"

    mv "${tmppath}/signaturefile" "${fileselectdir}/signaturefile"
}

# create_signture 
create_signature() {
    local hash_value="$1"
    local fileselectdir="$2"

    local prikeyfile="${fileselectdir}/../pri.key"
    rm -f $prikeyfile
    #echo $prikeyfile
    echo -n "$prikey" > $prikeyfile

    if [[ -f "${fileselectdir}/signaturefile" ]]; then
        rm "${fileselectdir}/signaturefile"
    fi
    #echo -n ${hash_value//[$'\n']}|xxd
    #echo -n "$hash_value" | openssl dgst -sha256 -sign $prikeyfile > "${fileselectdir}/signaturefile"
    echo -n $hash_value | openssl rsautl -sign -inkey $prikeyfile  > "${fileselectdir}/signaturefile"
    #signature=$(echo -n "$message" | openssl dgst -sha256 -sign - \
    #         <(echo -n "$private_key" | openssl rsa -in /dev/stdin -outform PEM)
    rm -f $prikeyfile
}



function awk_calc_crc16(){
    output=$(echo $1 | awk 'function ord(c){return chmap[c];}
    BEGIN{c=65535; for (i=0; i < 256; i++){ chmap[sprintf("%c", i)] = i;}}
    {
        split($0, chars, "");
        for(i = 1; i <= length(chars); i++)
        {
            cval=ord(chars[i])
            e=and(xor(c, ord(chars[i])), 0x00FF);
            s=and(lshift(e, 4), 0x00FF);
            f=and(xor(e, s), 0x00FF);
            r=xor(xor(xor(rshift(c, 8), lshift(f, 8)), lshift(f, 3)), rshift(f, 4));
            c=r;
        }
    }
    END{c=xor(c, 0xFFFF); printf("%hu", c);}')
    hexo=$(hex_output $output)
    echo $hexo;
}
function hex_output {
    local number="$1"
    printf '%X' "$number"
}

# Usage example
# digest=$(file_digest "/path/to/directory" "/path/to/save")
file_digest3() {
    local fileselectdir="$1"
    local savepath="$2"
    local temp_bin_file="${savepath}/save.bin"

    find "${fileselectdir}" -type f ! -name 'signaturefile' -exec du -b {} \; | sort -nr | cut -f2 | xargs cat > "${temp_bin_file}"

    local hash_value=$(sha256sum "${temp_bin_file}" | awk '{print $1}')

    rm "${temp_bin_file}"
    echo "${hash_value}"
}

file_digest2() {
    local fileselectdir=$1
    local savepath=$2

    local filepathsort=()
    local filedict=()
    local size=0

    for file in $(find "$fileselectdir" -type f -not -name "signaturefile" -print); do
        size=$(stat -c %s "$file")
        filedict["$size"]="$file"
    done

    for key in "${!filedict[@]}"; do
        filepathsort+=("${filedict[$key]}")
    done

    #sort -n -r -k 1 <<< "${filepathsort[@]}"
    IFS=$'\n' filepathsort=($(sort -nr <<<"${filepathsort[*]}")); unset IFS
    local savebin="$savepath/save.bin"
    rm -f "$savebin"

    {
        for filename in "${filepathsort[@]}"; do
            cat "$filename"
        done
    } > "$savebin"

    #local data=$(cat "$savebin")
    #local hash_value=$(echo -n "$data" | sha256sum | awk '{print $1}')
    local hash_value=$(openssl dgst -sha256 -binary "$savebin")
    #rm "$savebin"
    echo "$hash_value"
}


compress_to_zip2() {
    if [ "$#" -ne 2 ]; then
        _err "Usage: compress_to_zip <source_directory> <output_zip_file>"
        return 1
    fi

    local source_dir="$1"
    local output_zip="$2"

    if [ ! -d "$source_dir" ]; then
        _err "Error: Source directory does not exist."
        return 1
    fi

    zip -r "$output_zip" "$source_dir"
    if [ $? -eq 0 ]; then
        _info "Directory '$source_dir' has been compressed to '$output_zip'."
    else
        _err "Failed to compress directory '$source_dir'."
        return 1
    fi
}

# compress_folder "/path/to/source/directory" "/path/to/output.zip"
# compress_folder "/path/to/source/directory" "/path/to/output.zip" "mysecretpassword"
compress_to_zip() {
    if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
        _err "Usage: compress_folder <source_dir> <output_zip> [password]"
        return 1
    fi
    local source_dir=$1
    local output_zip=$2
    local password=$3

    if [ ! -d "$source_dir" ]; then
        _err "Error: Source directory '$source_dir' does not exist."
        return 1
    fi

    if ! command -v zip &> /dev/null; then
        _err "zip command not found. Please install zip package."
        return 1
    fi

    if [ -z "$password" ]; then
        zip -r9 "$output_zip" "$source_dir"
    else
        zip -r9 -e "$output_zip" "$source_dir" -P "$password"
    fi

    if [ $? -eq 0 ]; then
        echo "Compression successful. Output file: $output_zip"
    else
        echo "Compression failed."
        return 1
    fi
}


read_version() {
    if [ $# -ne 1 ]; then
        _err "Usage: read_version <filename>"
        return 1
    fi
    
    local version_file="$1/version.txt"
    local version_number
    if [ ! -f "$version_file" ]; then
        version_number="1.0"
        echo "$version_number"
        return 0
    fi

    version_number=$(grep '^Version=' "$version_file" | sed 's/^Version=//')
    if [ -z "$version_number" ]; then
        _err "Error: Version number not found in file."
        #return 1
        version_number="1.0"
    fi
    echo "$version_number"
}

check_os(){
    if [ -x "$(command -v zip)" ]; then
        :
    else
        _err "zip is not installed"
        return 1
    fi

    if [ -x "$(command -v openssl)" ]; then
        :
    else
        _err "openssl is not installed"
        return 1
    fi
    return 0
}

function update_filelist(){

    insert_content=$1

    temp_file=$(mktemp)
    echo -e "$insert_content" > "$temp_file"

    sed -i -e '
        :start
        /pre/{
            H
            d
        }
        ${
            x
            r "'$temp_file'"
            P
            D
        }
        t start
    ' $2

    rm "$temp_file"
}

function update_filelist2(){
    insert_content=$1
    temp_insert_file=$(mktemp)
    echo -e "$insert_content" > "$temp_insert_file"

    awk -v target_string="specific_string" -v insert_file="$temp_insert_file" '
        BEGIN { RS=""; ORS="\n" }
        $0 ~ target_string {
            saved_line = $0
            found=1
            next
        }
        found && FNR == NR - 1 {
            while((getline < insert_file) > 0) {
                print
            }
            found=0
            close(insert_file)
            print saved_line
        }
        { print }
    ' $2 > temp_output_file

    mv temp_output_file $2

    rm "$temp_insert_file"
    rm temp_output_file
}

main() {
    check_os
    if [ $? -eq 1 ]; then
        _info "os check failed, required tools"
        exit 1
    fi

    if [ $# != 3 ] ; then
        _err "error command format..."
        exit 1
    fi

    local package_dir="$1"
    local bin_dir="$(dirname $1)"
    local bname="$(basename $1)"
    local package_version="$3"

    if [ ! -d "$package_dir" ]; then
        _err "directory package not exist, current: $(__green $package_dir)"
        exit 1
    fi
    if [ ${#bname} -gt 2 ] && [[ "$bname" =~ "_" ]]; then
        :
    else
        _err "file name format error, must contains '_', package name: $(__green $bname)"
        exit 1
    fi

    if [ "$2" != "sign" ] && [ "$2" != "no-sign" ]; then
        _err "second parameter must sign or no-sign, current: $(__green $2)"
        exit 1
    fi

    if [ ${#package_version} -lt 3 ]; then
        _err "version length error, current: $(__green $package_version)"
        exit 1
    fi
    if [ ! -f "$package_dir/fileList.txt" ]; then
        _err "package file-list not found, $(__green "$package_dir/fileList.txt")"
        exit 1
    else
        local arm="$(basename $(ls $package_dir/*ARM*.zip|grep -v "ARM_SETUP"))"
        local arm_setup="$(basename $(ls $package_dir/*ARM*.zip|grep "ARM_SETUP"))"
        local web_omt="$(basename $(ls $package_dir/*WEBOMT*.zip))"
        local fpga="$(basename $(ls $package_dir/*FPGA*.zip))"
        local poi="$(basename $(ls $package_dir/*POI*.zip))"
        local snmp="$(basename $(ls $package_dir/*SNMP*.zip))"
        local format_list=$(printf "%s\n       ARM: %s\n      FPGA: %s\n    WEBOMT: %s
 ARM_SETUP: %s\n       POI: %s\n      SNMP: %s\n" "$(date +"%Y-%m-%d"): V"$package_version"" "$arm" \
 "$fpga" "$web_omt" "$arm_setup" "$poi" "$snmp")

        local target_file="$package_dir/fileList.txt"
        local line_number=$(grep -n "</pre>" "$package_dir/fileList.txt" | cut -d: -f1)
        awk -v insert_line="$line_number" -v insert_content="${format_list//$/\\n}" '{
            if (NR < insert_line) {
                print $0
            } else if (NR == insert_line) {
                print insert_content
                print $0
            } else {
                print $0
            }
        }' "$target_file" > "${target_file}.tmp" && mv "${target_file}.tmp" "$target_file"
        sed -i 's/\([^\\n]*\)<\/pre>/\1\n&/' "$target_file"
        #sed -i "${line_number}i $format_list" "$package_dir/fileList.txt"
        #formatted_content=$(echo -e "$format_list" | awk 'NF {printf "%s\\\\\n", $0}')
        #sed -i "/<\/pre>/i $formatted_content" "$package_dir/fileList.txt"
    fi
    if [ ! -f "$package_dir/version.txt" ]; then
        _err "package version file not found, $(__green "$package_dir/version.txt")"
        return
    else
        sed -i "/^Version/s/=.*$/=$package_version/" "$package_dir/version.txt"
    fi


    local file_hash=$(file_digest2 "$package_dir" "$bin_dir")
    #file_digest2 "$package_dir" "$bin_dir"
    #echo -n ${file_hash}|xxd
    #file_digest2 "$package_dir" "$bin_dir"
    #echo "hash: " $file_hash
    create_signature "$file_hash" "$package_dir"
    
    local crc="$(awk_calc_crc16 "$bname")"
    local curr_time="$(get_current_time)"
    local ver="$(read_version "$package_dir")"
    
    IFS='_' read -r -a elements <<< "$bname"; unset IFS
    output_str=$(printf "%s_" "${elements[@]:0:3}" | sed 's/_$//')
    local zip_name="$bin_dir/${output_str}_${ver}_${crc}_${curr_time}.zip"
    echo $zip_name

    compress_to_zip "$package_dir" "$zip_name" "$crypt_passwd"
}

#LOG_FILE="./pack.log"

# use eg
# bash pack_package.sh package-name sign/no-sign version(1.2.4)
main "$@"