<script>
  import { each } from "svelte/internal";

  let images = [];
  window.backend.Api.ListImages().then((result) => {
    images = result.map((image) => {
      return { name: image.RepoTags, status: "inactive" };
    });

    console.log(images);
  });

  function runContainer(imageName) {
    window.backend.Api.RunContainer(imageName[0]).then((result) => {
      images.forEach((image, index, _) => {
        if (image.name[0] == imageName[0]) {
          images[index] = { ...images[index], status: "running" };
        }
      });
    });
  }
</script>

<main>
  <table>
    <tr>
      <th>Image</th>
      <th>Action</th>
    </tr>
    {#each images as image}
      <tr>
        <td>{image.name}</td>
        <td>
          <button on:click={runContainer(image.name)}> Run </button>
        </td>
      </tr>
    {/each}
  </table>
</main>

<style></style>
