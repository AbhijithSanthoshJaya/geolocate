# Geolocate

This module is suppose to serve as a minimalistic client implementation
allowing M2M services and APIs to use Google Maps and Google Places APIs.
Reasons being:
Google did not update their official Go client to work with the New Places API. So API keys generated after March 1,2025 cannot use the official client library's implementation of places as it still incorrectly calls the old endpoint [here](https://github.com/googlemaps/google-maps-services-go/blob/eb7d4f974fd0540ee9be785e2bbae103d293562e/places.go#L35)
This would give API-KEY not authorised to error for these requests even if you enabled the Key to be used in the Google Cloud Console for the Old Places API explicitely. Not being able to use the Places API from Go client means you cannot do NearbySearch or TextSearch

## Local Setup

1. Clone the Repo
2. Ensure you have the latest Go installed and all PATH variables set
3. Add your API Key from Google Cloud Console. Follow instructions [here](https://developers.google.com/maps/documentation/javascript/get-api-key)
4. Make a .env file inside ./geo directory. Add `API_KEY=<replace with your API Key string>`. Save the file
5. Run the test functions in files with '\_test' to see the Google Maps and Places API responses

## Future Work

Built a Rest API that will help query and export these results
Build a Cache to look up last searched location results following Google's caching restrictions(30 days max)
