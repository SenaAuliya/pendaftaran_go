// src/types/app.d.ts
declare global {
	namespace App {
		// Custom error type for your application
		interface Error {
			message: string;
			statusCode: number;
		}

		// Custom data available in endpoints and hooks
		interface Locals {
			user: {
				id: string;
				name: string;
			};
		}

		// Custom page data for your pages
		interface PageData {
			title: string;
			description?: string;
		}

		// Custom state for pages
		interface PageState {
			previousPage: string;
		}

		// Custom platform data
		interface Platform {
			environment: 'production' | 'development';
		}
	}
}

export {};
