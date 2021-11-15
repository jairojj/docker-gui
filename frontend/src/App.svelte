<script>
	import Containers from "./components/Containers.svelte";
	import ImagesList from "./components/ImagesList.svelte";
	import Modal from "svelte-simple-modal";
</script>

<main>
	<div class="App">
		<header class="App-header">
			<h1>Images</h1>
			<ImagesList />
			<hr />
			<Modal>
				<Containers />
			</Modal>
		</header>
	</div>
</main>
