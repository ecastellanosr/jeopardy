document.addEventListener("click", function (event) {
    // Run your function
    if (event.target && event.target.id ==="first-button"){
        startcalculations();
    }
})

function startcalculations() {
    console.log("its working")
    let gridItem;
    const gridContainer = document.getElementById('cards-subgrid');
    const anycard = document.getElementById('card1');
    const dimensions = anycard.getBoundingClientRect();
    const cardWidth = dimensions.width
    const cardHeight = dimensions.height
    console.log("cardheight:",cardHeight,"cardWidth:",cardWidth)

    // Calculate left and bottom
    const GridWidth = gridContainer.offsetWidth; 
    const GridHeight = gridContainer.offsetHeight; 

    const WidthRatio = cardWidth/GridWidth
    const HeightRatio = cardHeight/GridHeight

    // Get the computed styles of the grid item

    gridContainer.addEventListener("click", function (event) {
        
    const target = event.target; // Typecast event.target to HTMLElement
    if (target.id && target.id.startsWith("card")) {
        gridItem = target; // Assign the clicked element
        const qanimation = document.querySelector(".Q-animation-holder");
        const gridIndex = Array.from(gridContainer.children).indexOf(gridItem);

        const columns = 6; // Number of columns in the grid
        const rows = 5
        let columnStart = (Math.floor((gridIndex+1)/ (columns-(0.9)))+1); // 1-based index
        
        let rowStart = ((gridIndex+1) % rows);
        if (rowStart === 0){
            rowStart = 5
        }

        const columnWidth = gridContainer.offsetWidth/6; // Assuming 6 columns
        const columnHeight = gridContainer.offsetHeight/6; // Assuming 6 rows

        rem = 16 // 1 rem is 16px
        // for x to be 0, we need to subtract half of a rem, 2 full rem because there are two columns, half of a width and 2 full widths
        const x = (-((rem/2)+(rem*((Math.ceil(columns/2))-1))+(cardWidth/2)+(cardWidth*((Math.ceil(columns/2))-1)))) + ((columnStart-1)*(rem+cardWidth)); 
        // for y to be 0 we add 2 rem and 2 column heights, to move columns
        const y = ((rem*((Math.ceil(columns/2))-1))+(cardHeight*((Math.ceil(columns/2))-1))) - ((rowStart-1)*(rem+cardHeight)); 
        qanimation.style.setProperty("--x",`${x}px`); 
        qanimation.style.setProperty("--y",`${y}px`);
        qanimation.style.setProperty("--width",`${GridWidth}px`);
        qanimation.style.setProperty("--height",`${GridHeight}px`);
        qanimation.style.setProperty("--Ratiox",`${WidthRatio}`);
        qanimation.style.setProperty("--Ratioy",`${HeightRatio}`);
        console.log("gridIndex",gridIndex,"columnWidth",columnWidth,"columnHeight",columnHeight,"columnStart",columnStart,"rowStart",rowStart,"x:",x,"y:",y,"Full Width:",GridWidth,"Full Height",GridHeight,"Width Ratio:",WidthRatio,"Heigth Ratio:",HeightRatio)

        // Trigger reflow to allow the browser to recognize the current state
        qanimation.offsetHeight; // This forces a reflow (do nothing with the result)

        if (qanimation) {
            
            qanimation.classList.add("Q-animation");
        }
    }
    });
}
