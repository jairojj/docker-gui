<script>
  import { each } from "svelte/internal";
  import { getContext } from "svelte";
  import Container from "./Container.svelte";

  const { open } = getContext("simple-modal");

  let containers = [];
  wails.Events.On("containerUpdate", (data) => (containers = data));

  function stopContainer(containerId) {
    window.backend.Api.StopContainer(containerId).then((result) => {});
  }
</script>

<main>
  <h1>Containers</h1>
  <table>
    <tr>
      <th>Name</th>
      <th>Image</th>
      <th>Status</th>
      <th>ID</th>
      <th>Action</th>
    </tr>
    {#each containers as container}
      <tr>
        <td>{container.Names}</td>
        <td>{container.Image}</td>
        <td>{container.Status}</td>
        <td>{container.Id}</td>
        <td>
          <button on:click={stopContainer(container.Id)}> Stop </button>
          <button
            on:click={open(
              Container,
              { container: container },
              { styleWindow: { width: "57rem" } }
            )}
          >
            Show Logs
          </button>
        </td>
      </tr>
    {/each}
  </table>
</main>

<style></style>
