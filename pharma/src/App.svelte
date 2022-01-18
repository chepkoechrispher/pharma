<script>
 import { onMount } from "svelte";
 import Edit from './Edit.svelte'
 import Delete from './Delete.svelte'

let myDrugs;
let id;
  onMount(async()=>{
     const response = await fetch("http://localhost:8000/drugs")
      const todo = await response.json();
      myDrugs = todo
      console.log(myDrugs)
  });
</script>

<div>
  {#if myDrugs}
  <h1>All Drugs</h1>
 
    <ul>
      {#each myDrugs as drug} 

      <p><i>id</i>  {drug.id} <i>name </i> {drug.name}  <i>name </i> {drug.stock}</p>
    
  {/each}
    </ul>
  {:else}
    <p>loading.....</p>
  {/if}


  <Edit/>
  <Delete/>
</div>



