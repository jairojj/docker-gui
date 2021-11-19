<script>
	export let container;
	let command;

	window.backend.Api.AttachToContainerShell(container.Id);

	function sendCommand() {
		if (command && command !== "" && command !== undefined) {
			console.log(command);
			wails.Events.Emit("container:attach:input", command);
		}
	}
</script>

<h1>{container.Names}</h1>

<input type="text" bind:value={command} />
<button on:click={sendCommand}> Send </button>

<style>
</style>
