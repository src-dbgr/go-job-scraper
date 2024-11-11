#!/bin/bash

output_file="output.txt"
preview_file="preview.txt"
max_file_size=$((1024 * 1024))  # 1 MB in Bytes
file_count=0
processed_count=0

# Arrays für die Ausschlussmuster
declare -a exclude_patterns
declare -a exclude_dirs

# Funktion zur Formatierung der Dateigröße
format_size() {
    local size=$1
    if ((size < 1024)); then
        echo "${size}B"
    elif ((size < 1048576)); then
        echo "$((size / 1024))KB"
    else
        echo "$((size / 1048576))MB"
    fi
}

# Funktion zur Überprüfung, ob eine Datei binär ist
is_binary() {
    if [[ "$(file -b --mime-encoding "$1")" == "binary" ]]; then
        return 0
    else
        return 1
    fi
}

# Funktion zur Überprüfung, ob eine Datei ausgeschlossen werden soll
should_exclude() {
    local file="$1"
    
    # Prüfe auf ausgeschlossene Verzeichnisse
    for dir in "${exclude_dirs[@]}"; do
        if [[ "$file" == *"/$dir/"* || "$file" == *"/$dir" ]]; then
            return 0
        fi
    done
    
    # Prüfe die Datei-Ausschlussmuster
    for pattern in "${exclude_patterns[@]}"; do
        if [[ "$file" == *"$pattern"* ]]; then
            return 0
        fi
    done
    return 1
}

# Funktion zur Erstellung einer Vorschau
create_preview() {
    local dir="$1"
    echo "Erstelle Vorschau für Verzeichnis: $dir" >> "$preview_file"
    while IFS= read -r -d '' file; do
        if [[ -f "$file" && "$file" != "./$output_file" && "$file" != "./$preview_file" ]]; then
            # Überprüfe Ausschlusskriterien
            if ! should_exclude "$file"; then
                local size=$(wc -c < "$file")
                if (( size <= max_file_size )) && ! is_binary "$file"; then
                    echo "Datei: $file ($(format_size $size))" >> "$preview_file"
                    ((file_count++))
                fi
            fi
        fi
    done < <(find "$dir" -type f -print0)
}

# Funktion zum Ausgeben des Inhalts einer Datei mit Pfad und Trennstrich
print_file_content() {
    local file="$1"
    local size=$(wc -c < "$file")
    echo "Verarbeite Datei: $file ($(format_size $size))"
    {
        echo "Datei: $file"
        echo "----------------------------------------"
        cat "$file"
        echo -e "\n----------------------------------------"
        echo ""
    } >> "$output_file"
    ((processed_count++))
    echo "Fortschritt: $processed_count / $file_count Dateien verarbeitet"
}

# Funktion zum Verarbeiten aller Dateien
process_files() {
    while IFS= read -r -d '' file; do
        if [[ -f "$file" && "$file" != "./$output_file" && "$file" != "./$preview_file" ]]; then
            # Überprüfe Ausschlusskriterien
            if ! should_exclude "$file"; then
                local size=$(wc -c < "$file")
                if (( size <= max_file_size )) && ! is_binary "$file"; then
                    print_file_content "$file"
                fi
            fi
        fi
    done < <(find . -type f -print0)
}

# Hilfe-Funktion
show_help() {
    echo "Verwendung: $0 [-e PATTERN] [-d DIR] ..."
    echo "  -e PATTERN    Dateimuster zum Ausschließen (kann mehrfach verwendet werden)"
    echo "  -d DIR        Verzeichnis zum Ausschließen (kann mehrfach verwendet werden)"
    echo "  -h           Diese Hilfe anzeigen"
    echo ""
    echo "Beispiele:"
    echo "  $0 -e \".log\" -e \"go.sum\" -d \".git\" -d \"node_modules\""
    echo "  $0 -d \".git\" -d \"vendor\" -e \".exe\""
    exit 1
}

# Parameter verarbeiten
while getopts "e:d:h" opt; do
    case $opt in
        e)
            exclude_patterns+=("$OPTARG")
            ;;
        d)
            exclude_dirs+=("$OPTARG")
            ;;
        h)
            show_help
            ;;
        \?)
            show_help
            ;;
    esac
done

# Hauptprogramm
echo "Erstelle Vorschau..." > "$preview_file"
create_preview "."
echo "Vorschau erstellt. $file_count Dateien gefunden."
echo "Ausgeschlossene Verzeichnisse: ${exclude_dirs[*]}"
echo "Ausgeschlossene Dateimuster: ${exclude_patterns[*]}"
echo "Möchten Sie die Vorschau anzeigen? (j/n)"
read -r show_preview
if [[ "$show_preview" == "j" ]]; then
    cat "$preview_file"
fi

echo "Möchten Sie fortfahren und alle Dateien verarbeiten? (j/n)"
read -r proceed
if [[ "$proceed" == "j" ]]; then
    : > "$output_file"  # Leert die Ausgabedatei
    echo "Starte Verarbeitung..."
    process_files
    echo "Verarbeitung abgeschlossen. Ergebnisse wurden in $output_file gespeichert."
else
    echo "Vorgang abgebrochen."
fi

# Aufräumen
rm "$preview_file"