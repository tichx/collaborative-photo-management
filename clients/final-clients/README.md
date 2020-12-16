# 441 Clients

Checkout the branch that has to do with your assignment and navigate to the build folder then run `live-server`.

# Auth branch

This will give you experience with going through the bundling process as part of deploying. The steps should be very simple compared to the previous time you deployed the client. This time however, React's default build script outputs the filies to a build folder. Use this information to your advantage.

## Setup

1. `npm install`. Any version of node _should_ work, but if something ends up failing to install it could be your version of node. I'm running node v12.3.0.
2. Modify the URL within `src/Constants/APIEndpoints/APIEndpoints.js` to match your API endpoint.
   1. `testbase` should be the API url when you run locally, so change the port number as needed.
   2. `base` is the API url hosted through your domain.
3. `npm run build` and the output files from build should be able to be copied over really easily.
4. The `index.html` file within the build folder is what React builds into. React produces a single page application unless you do some scripting.
5. `npm run start` for running on `localhost:3000`. You can also do `HTTPS=true npm run start` to run the app in HTTPS locally, I believe.
