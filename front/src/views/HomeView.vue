<script setup>
import { ref, watch } from "vue";
import axios from "axios";

const searchQuery = ref("");
const mapboxSearchResults = ref(null);
const showBody = ref(null);

const currentPage = ref(1);
const pagesQuantity = ref(5);

watch(currentPage, () => getSearchResults());

const getSearchResults = async () => {
    if (searchQuery.value !== "") {
        const response = await axios.post(
            `http://localhost:3000/search/${searchQuery.value}&page=${currentPage.value}`
        );
        mapboxSearchResults.value = response.data; // array
        pagesQuantity.value = mapboxSearchResults.pagesQuantity;
        console.log(mapboxSearchResults.value);
        return;
    }
    mapboxSearchResults.value = null;
};
</script>

<template>
    <main class="container ">
        <div class="py-8 px-4 mx-auto max-w-screen-xl">
             <div class="pt-4 mb-8 relative">
                <input
                    type="text"
                    v-model="searchQuery"
                    @input="getSearchResults"
                    placeholder="Search emails"
                    class="text-white py-2 px-1 w-full bg-transparent border-b focus:outline-none focus:ring ring-blue-950"
                 />
            </div>
            <div class="flex flex-1">
            <div class="relative overflow-x-auto shadow-md sm:rounded-lg text-white ">
            <table 
            v-if="Array.isArray(mapboxSearchResults)"
            class="w-full text-sm text-left text-gray-500 dark:text-gray-400"
            >
            
                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                        <th  scope="col" class="px-6 py-3">Subject</th>
                        <th  scope="col" class="px-6 py-3">From</th>
                        <th  scope="col" class="px-6 py-3">To</th>
                        <th  scope="col" class="px-6 py-3">View</th>
                    </tr>
                </thead>
                <tbody>
                    <tr
                    class="bg-white border-b dark:bg-gray-900 dark:border-gray-700"
                        v-for="(item, index) in mapboxSearchResults"
                        :key="index"
                    >
                        <td class="px-6 py-4">{{ item._source.subject }}</td>
                        <td class="px-6 py-4">{{ item._source.from }}</td>
                        <td class="px-6 py-4">{{ item._source.to }}</td>
                        <td class="px-6 py-4">
                            <button @click="showBody = item._source.Body">
                                <i class="fa-solid fa-eye text-2xl"></i>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        </div>


        
        <span class=" text-xs leading-none px-4  dark:text-blue-200 absolute "> 
            <div class="bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg md:p-12 fixed  left-1/2 top-40 bottom-5 right-5 overflow-y-auto max-h-screen"
            >
                
                        <div  v-if="!!showBody">
                          {{ showBody }}
                         </div>
     
                      
            </div>  
        </span> 
    </div>
    </main>
</template>
