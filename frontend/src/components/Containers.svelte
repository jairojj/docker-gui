<script>
  import { each } from "svelte/internal";

  let containers = [];
  wails.Events.On("containerUpdate", (data) => (containers = data));
</script>

<main>
  <h1>Containers</h1>
  <ul>
    {#each containers as container}
      <li>
        {container.Names} - {container.Image} - {container.Status}
      </li>
    {/each}
  </ul>
</main>

<style></style>
