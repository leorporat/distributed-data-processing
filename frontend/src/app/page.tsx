"use client";

import { useState, FormEvent } from "react";
import ResultsDisplay, { SummaryResult } from "@/components/ResultsDisplay";

// Define TypeScript interfaces for our form state
interface RedditSearchForm {
  query: string;
  subreddit: string;
  limit: number;
  isLoading: boolean;
}

// Popular subreddits for suggestions
const POPULAR_SUBREDDITS = [
  "AskReddit",
  "explainlikeimfive",
  "todayilearned",
  "science",
  "technology",
  "programming",
  "datascience",
];

export default function Home() {
  // Initialize form state with default values
  const [formState, setFormState] = useState<RedditSearchForm>({
    query: "",
    subreddit: "",
    limit: 10,
    isLoading: false,
  });

  // State for search results
  const [searchError, setSearchError] = useState<string | undefined>();
  const [searchResult, setSearchResult] = useState<SummaryResult | undefined>();

  // Handle form submission
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    // Reset previous results and errors
    setSearchError(undefined);
    setFormState({ ...formState, isLoading: true });
    
    try {
      // Log the search parameters
      console.log("Searching for:", formState);
      
      // Import the Reddit service
      const { searchReddit } = await import("@/services/redditService");
      
      try {
        // Make the API call to our backend
        const result = await searchReddit(
          formState.subreddit,
          formState.query,
          formState.limit
        );
        setSearchResult(result);
      } catch (error) {
        console.error("Error from Reddit service:", error);
        setSearchError(`Error fetching data: ${error instanceof Error ? error.message : String(error)}`);
      } finally {
        setFormState({ ...formState, isLoading: false });
      }
    } catch (err) {
      setSearchError(`An error occurred: ${err}`);
      setFormState({ ...formState, isLoading: false });
    }
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 bg-gradient-to-b from-white to-gray-100 dark:from-gray-900 dark:to-black">
      <main className="w-full max-w-2xl mx-auto">
        {/* Header */}
        <div className="text-center mb-10">
          <h1 className="text-4xl font-bold mb-2 text-gray-800 dark:text-white">
            Reddit Summarizer
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-300">
            Get concise summaries of Reddit discussions
          </p>
        </div>

        {/* Search Form */}
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 mb-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Search Query Input */}
            <div>
              <label 
                htmlFor="query" 
                className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1"
              >
                Search Query
              </label>
              <input 
                type="text"
                id="query"
                placeholder="What do you want to search for on Reddit?"
                className="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={formState.query}
                onChange={(e) => setFormState({ ...formState, query: e.target.value })}
                required
              />
            </div>

            {/* Subreddit Input */}
            <div>
              <label 
                htmlFor="subreddit" 
                className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1"
              >
                Subreddit
              </label>
              <input 
                type="text"
                id="subreddit"
                placeholder="e.g., AskReddit, programming, science"
                className="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={formState.subreddit}
                onChange={(e) => setFormState({ ...formState, subreddit: e.target.value })}
                required
              />
              <div className="mt-2 flex flex-wrap gap-2">
                {POPULAR_SUBREDDITS.map((subreddit) => (
                  <button
                    key={subreddit}
                    type="button"
                    className="px-2 py-1 text-xs rounded-full bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-blue-100 dark:hover:bg-blue-900 transition-colors"
                    onClick={() => setFormState({ ...formState, subreddit })}
                  >
                    r/{subreddit}
                  </button>
                ))}
              </div>
            </div>

            {/* Limit Input */}
            <div>
              <label 
                htmlFor="limit" 
                className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1"
              >
                Result Limit
              </label>
              <select
                id="limit"
                className="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={formState.limit}
                onChange={(e) => setFormState({ ...formState, limit: parseInt(e.target.value) })}
              >
                <option value="5">5 posts</option>
                <option value="10">10 posts</option>
                <option value="25">25 posts</option>
                <option value="50">50 posts</option>
              </select>
            </div>

            {/* Submit Button */}
            <button 
              type="submit"
              className="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={formState.isLoading || !formState.query || !formState.subreddit}
            >
              {formState.isLoading ? 'Searching...' : 'Search Reddit & Summarize'}
            </button>
          </form>
        </div>

        {/* Results Display Component */}
        <ResultsDisplay 
          isLoading={formState.isLoading}
          error={searchError}
          result={searchResult}
          originalQuery={formState.query}
          subreddit={formState.subreddit}
        />
      </main>

      {/* Footer */}
      <footer className="mt-auto pt-8 text-center text-gray-500 dark:text-gray-400 text-sm">
        <p>Reddit Summarizer - Get insights without the noise</p>
      </footer>
    </div>
  );
}
