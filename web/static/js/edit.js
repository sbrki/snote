const $editor = document.getElementById("editor");
const $vimCheckbox = document.getElementById("vim-checkbox");
const $vimCheckBoxLabel = document.getElementById("vim-checkbox-label");
const $previewUrl = document.getElementById("preview-url");
const $saveButton = document.getElementById("save-button");
const $saveButtonAnimator = document.getElementById("save-button-animator");

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const editor =
	CodeMirror.fromTextArea(
		$editor,
		{
			lineNumbers: true,
			theme: "default",
			mode: "markdown",
			keyMap: "default",
			showCursorWhenSelecting: true
		},
	)
;
editor.setSize("100%", "100%");

// Vim mode switch
$vimCheckbox.addEventListener("change", () => { 
	if ($vimCheckbox.checked == true) {
		editor.setOption("keyMap", "vim");
		$vimCheckBoxLabel.style.color = "green";
	} else {
		editor.setOption("keyMap", "default");
		$vimCheckBoxLabel.style.color = null;
	}

	editor.focus();
});

const [ , currentNoteId ] = window.location.pathname.split("/");

// set the preview href
$previewUrl.href = `/${currentNoteId}`;

// Fetch JSON template for note
let noteJson;
(async () => {
	const response = await fetch(`/api/note/${currentNoteId}`);

	if (response.status !== 200) {
		if (response.status === 404) {
			alert("Note not found (404)");
		}
		else {
			alert("Error sending GET request to fetch note");
			alert(xhr.status);
		}
		return;
	}

	noteJson = await response.json();
	const { contents } = noteJson;

	$editor.value = contents;
	editor.getDoc().setValue(contents);
})();

async function saveNote() {
	$saveButton.style.color = "gray";

	noteJson.contents = editor.getDoc().getValue();

	const response = await fetch(
		`/api/note/${currentNoteId}`,
		{
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(noteJson),
		},
	);

	if (response.status !== 200) {
		$saveButton.style.color = "red";
		$saveButton.innerHTML = "ERROR SAVING";

		return;
	}

	$saveButton.style.color = "green";
	$saveButton.innerHTML = "saved!";

	await animateCSS("#save-button", "jello");

	$saveButton.style.color = null;
	$saveButton.innerHTML = "save";
}

// automatically call saveNote() on every 100 keystrokes
let keystrokeNum = 0;
editor.on("change", async () => {
	keystrokeNum += 1;

	if (keystrokeNum >= 100) {
		await saveNote();
		keystrokeNum = 0;

		return;
	}
});

// call saveNote() on ctrl+s (or CMD+s)
document.addEventListener("keydown", async (e) => {
	const isCtrlKeyPressed = e.metaKey || e.ctrlKey;
	const isSKeyPressed = e.code === "KeyS";

	if (isCtrlKeyPressed && isSKeyPressed) {
		e.preventDefault();
		await saveNote();
	}
}, false);

async function newNotePrompt() {
	const newNoteId = window.prompt("Enter ID of the new note:\n(the note will live at /<note ID>)");

	const noteExistsResponse = await fetch(`/api/note/${newNoteId}`);
	if (noteExistsResponse.status !== 404) {
		alert("note already exists!");
		return;
	}

	const prototypeResponse = await fetch("/api/note_prototype");
	if (prototypeResponse.status !== 200) {
		alert("error fetching note prototype")
		return;
	}

	const prototype = await prototypeResponse.json();
	prototype.id = newNoteId;
	prototype.contents = `# ${newNoteId}`;

	const createNoteResponse = await fetch(
		`/api/note/${newNoteId}`,
		{
			method: "PUT",
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(prototype),
		},
	);

	if (createNoteResponse.status === 201) {
		window.location.replace(`/${newNoteId}/edit`);
	} else {
		alert("Error creating new note (failed to PUT prototype)");
	}
}

async function deleteNotePrompt() {
	const confirmation = window.prompt("Type 'yes' to confirm deletion");

	if (confirmation !== "yes") {
		return;
	}

	const deleteResponse = await fetch(
		`/api/note/${currentNoteId}`,
		{
			method: "DELETE",
		},
	);

	if (deleteResponse.status === 200) {
		window.location.replace("/ls");
	} else {
		alert("Error deleting note.");
	}
}


// animate.css utility function taken from official docs
const animateCSS = (element, animation, prefix = 'animate__') => {
	// We create a Promise and return it
	return new Promise((resolve, reject) => {
		const animationName = `${prefix}${animation}`;
		const node = document.querySelector(element);

		node.classList.add(`${prefix}animated`, animationName);

		// When the animation ends, we clean the classes and resolve the Promise
		function handleAnimationEnd() {
			node.classList.remove(`${prefix}animated`, animationName);
			resolve('Animation ended');
		}

		node.addEventListener('animationend', handleAnimationEnd, {once: true});
	}
	);
}
