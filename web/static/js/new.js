async function newNotePrompt() {
	const newNoteId = window.prompt("Enter ID of the new note:\n(the note will live at /<note ID>)");
	let formData = new FormData();
	formData.append("suggested_id", newNoteId);
	const newNoteResponse = await fetch(`/api/note`,
		{
			method: "POST",
			body: formData,
		},
	);

	if (newNoteResponse.status === 409) {
		alert("Note already exists!");
		return;
	} else if (newNoteResponse.status !== 201) {
		alert("Error creating note!");
		alert(newNoteResponse.status);
		return;
	}
	window.location.replace(`/${newNoteId}/edit`);
}
