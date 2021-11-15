<script>
  export let container;
  let logs = null;

  window.backend.Api.ContainerLogs(container.Id).then((result) => {
    logs = result;
  });
</script>

<h1>{container.Names}</h1>

<pre>
  {logs}
</pre>

<style>
  pre {
    white-space: break-spaces;
  }
</style>
