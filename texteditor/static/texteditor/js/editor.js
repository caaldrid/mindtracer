let buffer = [];
let gapIndex = 0;
const allowedKeys = /^[a-zA-Z0-9 .,!?;:'`~"()\[\]{}<>@#$%^&*_+=\\/-]$/;
const caretHTML = "\<span class=\"caret\" id=\"caret\"\>\<\/span\>";

function clearSelection() {
    let selectedLines = htmx.findAll(".selected");
    let newLine = htmx.find(".new");
    let prevLine = null;
    selectedLines.forEach((line, _) => {
        if (!line.classList.contains("new")) {
            htmx.toggleClass(line, "selected"); // Remove the "selection" class from the previously selected line
        }
    });

    if (newLine != null) {
        htmx.toggleClass(newLine, "new"); // Remove the "new" class from the new line now that we have cleared the previous selection
    }
}

function clearCaret() {
    let line = htmx.find(".selected > .text")
    let caret = htmx.find(line, ".caret")

    if (caret != null) {
        htmx.remove(caret)
        buffer = Array.from(line.innerText);
    }
}

function moveCaret() {
    buffer.splice(gapIndex, 0, caretHTML);
    let newText = buffer.join('') || "";
    htmx.swap(".selected > .text", newText, { swapStyle: "innerHTML" })
}

function updateLineNumbers() {
    let lines = htmx.findAll(".line");
    lines.forEach((line, index) => {
        htmx.find(line, ".number").innerText = index + 1;
    });
}

function selectLine(event) {
    clearCaret();
    clearSelection();
    let lineClickedOn = event.currentTarget;
    htmx.toggleClass(lineClickedOn, "selected"); // Toggle the "selected" class the line we clicked on

    let lineText = htmx.find(lineClickedOn, ".text");
    buffer = Array.from(lineText.innerText);

    const selection = window.getSelection();
    if (!selection.rangeCount) return;

    let range = selection.getRangeAt(0);
    // Ensure the range is inside our div
    if (!lineText.contains(range.commonAncestorContainer)) return;

    let preCaretRange = range.cloneRange();
    preCaretRange.selectNodeContents(lineText);
    preCaretRange.setEnd(range.startContainer, range.startOffset);
    gapIndex = preCaretRange.toString().length;
    moveCaret();
}

function handleKeydown(event) {
    clearCaret();

    if (event.key.length === 1 && allowedKeys.test(event.key)) {
        buffer.splice(gapIndex, 0, event.key);
        gapIndex++;
    } else if (event.key === 'Backspace') {
        if (gapIndex == 0) {
            let selected = htmx.find(".selected");
            let remainingText = htmx.find(selected, ".text").innerText;

            let lineAbove = selected.previousElementSibling;
            if (lineAbove != null) {
                let lineAboveText = htmx.find(lineAbove, ".text").innerText
                htmx.remove(selected);
                htmx.toggleClass(lineAbove, "selected");

                // Append the text of the deleted line with whats in the line above 
                buffer = Array.from(lineAboveText + remainingText);
                gapIndex = buffer.length
            }
        } else {
            gapIndex--;
            buffer.splice(gapIndex, 1);
        }
    } else if (event.key === 'ArrowLeft' && gapIndex > 0) {
        gapIndex--;
    } else if (event.key === 'ArrowRight' && gapIndex < buffer.length) {
        gapIndex++;
    } else if (event.key === 'ArrowUp' || event.key === 'ArrowDown') {
        let selected = htmx.find(".selected");
        let selectedText = htmx.find(selected, ".text").innerText;

        let target = event.key === 'ArrowUp' ? selected.previousElementSibling : selected.nextElementSibling;
        if (target != null) {
            clearSelection();
            buffer = Array.from(htmx.find(target, ".text").innerText);
            gapIndex = buffer.length < selectedText.length ? buffer.length : gapIndex;
            htmx.toggleClass(target, "selected");
        }
    } else if (event.key === 'Enter') {
        event.preventDefault();
        // Set the text of currently selected line to everything before the gapIndex
        let newline = buffer.slice(0, gapIndex);
        htmx.find(".selected > .text").innerText = newline.join('') || "";

        buffer = buffer.slice(gapIndex, buffer.length);
        gapIndex = 0;

        htmx.swap(".selected", "<span class=\"line selected new\" hx-on:click=\"selectLine(event)\"><span class=\"number\">1</span><span class=\"text\"></span></span>", { swapStyle: "afterend" });
        clearSelection();
        updateLineNumbers();
    }

    moveCaret();
}

