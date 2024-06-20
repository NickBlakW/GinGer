window.addEventListener("load", () => {
	addTitle();
});

function addTitle() {
	const elemDiv = document.createElement("div");
	const title = document.createElement("h1");
	title.innerText = "GinGer API Viewer";

	elemDiv.appendChild(title);

	document.body.insertAdjacentElement("afterbegin", elemDiv);
}
