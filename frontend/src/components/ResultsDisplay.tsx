"use client";

import React from 'react';

// TypeScript interfaces for the Reddit API response data structure
export interface RedditComment {
  body: string;
}

export interface RedditPost {
  postTitle: string;
  comments: RedditComment[];
}

// Interface for the summarized content
export interface SummaryResult {
  originalQuery: string;
  subreddit: string;
  postTitle: string;
  summary: string;
  sentimentScore?: number; // Optional for future use
  topComments?: RedditComment[]; // Optional for showing selected comments
}

// Props for the ResultsDisplay component
interface ResultsDisplayProps {
  isLoading: boolean;
  error?: string;
  result?: SummaryResult;
  originalQuery: string;
  subreddit: string;
}

const ResultsDisplay: React.FC<ResultsDisplayProps> = ({
  isLoading,
  error,
  result,
  originalQuery,
  subreddit
}) => {
  // Loading state
  if (isLoading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 animate-pulse">
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-6"></div>
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"></div>
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-5/6 mb-2"></div>
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-4"></div>
        <div className="h-20 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border-l-4 border-red-500">
        <h3 className="text-lg font-semibold text-red-600 dark:text-red-400 mb-2">Error</h3>
        <p className="text-gray-700 dark:text-gray-300">{error}</p>
        <div className="mt-4 p-4 bg-gray-100 dark:bg-gray-700 rounded-lg">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            <span className="font-medium">Query:</span> {originalQuery}
          </p>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            <span className="font-medium">Subreddit:</span> r/{subreddit}
          </p>
        </div>
      </div>
    );
  }

  // No results yet
  if (!result) {
    return null;
  }

  // Results display
  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
      {/* Header with search info */}
      <div className="mb-6 pb-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex flex-col md:flex-row md:justify-between md:items-center">
          <h2 className="text-xl font-bold text-gray-800 dark:text-white mb-2 md:mb-0">
            Summary Results
          </h2>
          <div className="text-sm text-gray-600 dark:text-gray-400">
            <span className="font-medium">Query:</span> {result.originalQuery} | 
            <span className="font-medium ml-2">Subreddit:</span> r/{result.subreddit}
          </div>
        </div>
      </div>

      {/* Post title */}
      <div className="mb-4">
        <h3 className="text-lg font-medium text-gray-800 dark:text-white">{result.postTitle}</h3>
        <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
          From Reddit discussion
        </div>
      </div>

      {/* Summary */}
      <div className="mb-6 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <h4 className="text-md font-medium text-gray-800 dark:text-white mb-2">Summary</h4>
        <p className="text-gray-700 dark:text-gray-300 whitespace-pre-line">
          {result.summary}
        </p>
      </div>

      {/* Top comments section (if available) */}
      {result.topComments && result.topComments.length > 0 && (
        <div>
          <h4 className="text-md font-medium text-gray-800 dark:text-white mb-3">
            Notable Comments
          </h4>
          <div className="space-y-3">
            {result.topComments.map((comment, index) => (
              <div 
                key={index} 
                className="p-3 bg-gray-50 dark:bg-gray-700 rounded-lg border-l-2 border-blue-500"
              >
                <p className="text-gray-700 dark:text-gray-300">{comment.body}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Sentiment score (if available) */}
      {result.sentimentScore !== undefined && (
        <div className="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
          <h4 className="text-sm font-medium text-gray-800 dark:text-white mb-2">
            Discussion Sentiment
          </h4>
          <div className="flex items-center">
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
              <div 
                className={`h-2.5 rounded-full ${
                  result.sentimentScore > 0.6 
                    ? 'bg-green-500' 
                    : result.sentimentScore > 0.4 
                      ? 'bg-yellow-500' 
                      : 'bg-red-500'
                }`}
                style={{ width: `${result.sentimentScore * 100}%` }}
              ></div>
            </div>
            <span className="ml-2 text-sm text-gray-600 dark:text-gray-400">
              {Math.round(result.sentimentScore * 100)}%
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default ResultsDisplay;

