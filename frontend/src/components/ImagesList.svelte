<script>
  import { each } from "svelte/internal";

  let images = [];
  window.backend.Api.ListImages().then((result) => {
    images = result;
  });
</script>

<main>
  <ul>
    {#each images as image}
      <li>
        {image.RepoTags}
      </li>
    {/each}
  </ul>
</main>

<style></style>
