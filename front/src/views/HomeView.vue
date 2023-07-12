<script setup>
import { ref, watch } from "vue";
import axios from "axios";

const searchQuery = ref("");
const mapboxSearchResults = ref(null);
const showBody = ref(null);

const currentPage = ref(0);
const hasNextPage = ref(true)

watch(currentPage, () => getSearchResults());
watch(searchQuery, () => {
    if (currentPage.value === 0) {
        getSearchResults()
    } else {
        currentPage.value = 0
    }

    hasNextPage.value = true
});

const getSearchResults = async () => {
    if (searchQuery.value !== "") {
        const response = await axios.get(
            `http://localhost:3000/search/${searchQuery.value}?page=${currentPage.value}`
        );
        mapboxSearchResults.value = response.data; // array

        hasNextPage.value = mapboxSearchResults.value.length > 0

        return;
    }
    mapboxSearchResults.value = null;
};
</script>

<template>
    <main>
        <div class="py-8 px-4 mx-auto max-w-screen-xl flex justify-between">
            <div class="flex-1">

                <div class="pt-4 mb-8 relative">
                    <input type="text" v-model="searchQuery" placeholder="Search emails"
                        class="text-white py-2 px-1 w-full bg-transparent border-b focus:outline-none focus:ring ring-blue-950" />
                </div>
                <div class="flex flex-1">
                    <div class="relative overflow-x-auto shadow-md sm:rounded-lg text-white ">
                        <table v-if="Array.isArray(mapboxSearchResults)"
                            class="w-full text-sm text-left text-gray-500 dark:text-gray-400">

                            <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                                <tr>
                                    <th scope="col" class="px-6 py-3">Subject</th>
                                    <th scope="col" class="px-6 py-3">From</th>
                                    <th scope="col" class="px-6 py-3">To</th>
                                    <th scope="col" class="px-6 py-3">View</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr class="bg-white border-b dark:bg-gray-900 dark:border-gray-700"
                                    v-for="(item, index) in mapboxSearchResults" :key="index">
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

                <nav v-if="Array.isArray(mapboxSearchResults)" class="flex w-full justify-center mt-4" aria-label="Page navigation example">

                    <button v-show="currentPage > 0" @click="currentPage--"
                        class="flex items-center justify-center px-3 h-8 ml-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                        <span class="sr-only">Previous</span>
                        <svg class="w-2.5 h-2.5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                            viewBox="0 0 6 10">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M5 1 1 5l4 4" />
                        </svg>
                    </button>

                    <button v-show="hasNextPage" @click="currentPage++"
                        class="flex items-center justify-center px-3 h-8 leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                        <span class="sr-only">Next</span>
                        <svg class="w-2.5 h-2.5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                            viewBox="0 0 6 10">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="m1 9 4-4-4-4" />
                        </svg>
                    </button>
                </nav>

            </div>


            <div class="flex-1 text-xs leading-none px-4  dark:text-blue-200 break-words overflow-hidden ">
                <div class="bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg md:p-12 ">

                    <div v-if="!!showBody">
                        {{ showBody }}
                    </div>


                </div>
            </div>
        </div>

    </main>
</template>
