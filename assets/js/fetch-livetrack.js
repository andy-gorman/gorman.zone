import * as params from '@params';

window.addEventListener("load", async (event) => {
	const response = await fetch(
		`${params.apiURL}/garmin-live-track`,
		{ cache: 'no-store' }
	);
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
	const frameWrapper = document.getElementById("garmin-livetrack-wrapper");
	frameWrapper.innerHTML = `<h3>Andy is not riding right now. Check Back Later</h3>`;
});
