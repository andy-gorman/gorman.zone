import * as params from '@params';

window.addEventListener("load", async (event) => {
	const response = await fetch(`${params.apiURL}/garmin-live-track`);
	if (!response.ok) {
		console.error(`Response status: ${response.status}`);
		return;
	}
	const url = await response.text();
	if (url !== "") {
		const garminFrame = document.getElementById("garmin-livetrack");
		garminFrame.setAttribute("src", url);
		return;
	}
	// TODO: Render text when Andy isn't riding when server returns an empty string
});
