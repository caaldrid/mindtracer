let buffer = [];
let gapIndex = 0;
const allowedKeys = /^[a-zA-Z0-9 .,!?;:'`~"()\[\]{}<>@#$%^&*_+=\\/-]$/;
const caretHTML = "\<span class=\"caret\" id=\"caret\"\>\<\/span\>";

function clearCaret() {
    let line = htmx.find(".selected > .text");
    let caret = htmx.find(line, ".caret");

    if (caret != null) {
        htmx.remove(caret);
        buffer = Array.from(line.innerText);
    }
}

function moveCaret() {
    buffer.splice(gapIndex, 0, caretHTML);
    let newText = buffer.join('') || "";
    htmx.swap(".selected > .text", newText, { swapStyle: "innerHTML" });
}

function updateLineNumbers() {
    let lines = htmx.findAll(".line");
    lines.forEach((line, index) => {
        htmx.find(line, ".number").innerText = index + 1;
    });
}

function selectLine(event) {
    clearCaret();
    htmx.toggleClass(htmx.find(".selected"), "selected")
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
    if (event.ctrlKey) {
        return;
    }
    clearCaret();
    let key = event.key;
    let selected = htmx.find(".selected");

    switch (key) {
        case 'Backspace':
            if (gapIndex == 0) {
                let remainingText = htmx.find(selected, ".text").innerText;

                let lineAbove = selected.previousElementSibling;
                if (lineAbove != null) {
                    let lineAboveText = htmx.find(lineAbove, ".text").innerText;
                    htmx.remove(selected);
                    htmx.toggleClass(lineAbove, "selected");

                    // Append the text of the deleted line with whats in the line above 
                    buffer = Array.from(lineAboveText + remainingText);
                    gapIndex = buffer.length;
                }
            } else {
                gapIndex--;
                buffer.splice(gapIndex, 1);
            }
            break;
        case 'ArrowLeft':
            if (gapIndex > 0) {
                gapIndex--;
            }
            break;
        case 'ArrowRight':
            if (gapIndex < buffer.length) {
                gapIndex++;
            }
            break;
        case 'ArrowUp':
        case 'ArrowDown':
            let selectedText = htmx.find(selected, ".text").innerText;
            let target = event.key === 'ArrowUp' ? selected.previousElementSibling : selected.nextElementSibling;
            if (target != null) {
                htmx.toggleClass(htmx.find(".selected"), "selected")
                
                buffer = Array.from(htmx.find(target, ".text").innerText);
                gapIndex = buffer.length < selectedText.length ? buffer.length : gapIndex;
                htmx.toggleClass(target, "selected");
            }
            break;
        case 'Enter':
            event.preventDefault();
            // Split the text of currently selected line at the gapIndex
            let newline = buffer.slice(0, gapIndex);
            htmx.find(selected, ".text").innerText = newline.join('') || "";

            buffer = buffer.slice(gapIndex, buffer.length);
            gapIndex = 0;


            htmx.ajax("GET", "/new-line", {target: ".selected", swap: "afterend"}).then(() => {
                htmx.toggleClass(htmx.find(".selected"), "selected");
                updateLineNumbers();
                moveCaret();
            });
            return;
        default:
            if  (allowedKeys.test(key)) {
                buffer.splice(gapIndex, 0, event.key);
                gapIndex++;
            }
            break;
    }

    moveCaret();
}

