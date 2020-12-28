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
	editor.save = saveNote;
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

// file upload (via dropzone.js)
let dropzone = new Dropzone("#file-upload",
	{ 
		url: "/api/blob",
		maxFilesize: 2048, // 2GiB
		timeout: 300000, // 5 minutes
		clickable: false,
		init: function() {
			this.on("complete", function(f) {
				if (f.xhr.status == 201) {
					let loc = f.xhr.getResponseHeader("Location");
					let doc = editor.getDoc();
					let cur = doc.getCursor();
					let filename = f.name;
					let filesize = (f.size / (1000**2)).toFixed(2);
					doc.replaceRange(`![${filename} (${filesize}M)](${loc})\n`, cur);
					saveNote();
				} else {
					alert("file upload failed");
					alert(f.xhr.status);
				}
				this.removeFile(f);
			});
		},
	},
);

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
