<script>
 import { onMount } from "svelte";

let myDrugs;
let id;
let url = 'http://localhost:8000/drugs/'+''+id
  

  const handleEdit = (e) => {
    onMount(async()=>{
    const response = await fetch(url)
      const todo = await response.json();
      myDrugs = todo
      console.log(myDrugs)
  }); 
  }
</script>

<div class="Edit">
  
  <h1>Edit</h1>
  <h3>Enter ID to edit stock<h3>
  <input type=number bind:value={id}>
  <button on:click={handleEdit} value="Edit">Edit</button>
</div>
