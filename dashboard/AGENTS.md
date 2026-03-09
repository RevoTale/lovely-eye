- use only Typescript: `.tsx` in case you need to defined the component, otherwise `.ts` files
- defined only one component per file
- maximum file size should be 120 lines. Split it into separate modules/components that are naming by what they do provide.
- use `FunctionComponent` for the Next.js/React components.
- name of the component defined should match the file name.
- for React components use the `export default` keyword for exporting.

# Data loading behaviour
- each component should have the fallback skeleton which repeats the visual structure
    - for inline elements (with text) the height of single row should match the skeleton row height
- prefer `cache-and-network` policy for the Graphql data fetch. 
- to indicate that data is being loaded create a poper loader indicators along with displaying the stale data until new data will load
- pagination/navigation should not cause flickering. Elements height have the architecture to avoid the "layout jumping".
