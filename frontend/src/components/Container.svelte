<script>
  export let container;
  let logs = null;

  window.backend.Api.ListenForContainerLogs(container.Id);

  wails.Events.On("container:log:new", (data) => {
    logs = data;
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
