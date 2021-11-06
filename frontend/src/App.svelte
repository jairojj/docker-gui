<script>
	import Containers from "./components/Containers.svelte";
	import ImagesList from "./components/ImagesList.svelte";
</script>

<main>
	<div class="App">
		<header class="App-header">
			<h1>Images</h1>
			<ImagesList />
			<hr />
			<Containers />
		</header>
	</div>
</main>
