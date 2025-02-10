window.addEventListener("load", (event) => {
    startcalculations();
});


function startcalculations() {
    let gridItem;
    const gridContainer = document.getElementById('cards-subgrid');
    const anycard = document.getElementById('card1');
    //card dimensions
    const dimensions = anycard.getBoundingClientRect();
    const cardWidth = dimensions.width
    const cardHeight = dimensions.height
    console.log("cardheight:",cardHeight,"cardWidth:",cardWidth) //for troubleshooting 

    // Calculate Grid width and height to get position of columns and rows with card's dimensions
    const GridWidth = gridContainer.offsetWidth; 
    const GridHeight = gridContainer.offsetHeight; 

    const WidthRatio = cardWidth/GridWidth
    const HeightRatio = cardHeight/GridHeight

    // qanimation is the canvas for the card animation, the card itself isnt animating, we replace it with this canvas. 
    // its done like this because we had multiple things that needed to be done with the card, 
    // it needs to place another empty card to not move the cards position, send websocket to show the question after the animation and remove cover for question. 
    // adding the animation to it would add complexity to the card element

    let qanimation = document.querySelector(".Q-animation-holder");
    // event listener for a websocket send before its done as the card item sends a websocket message to fetch the question of the card.
    document.addEventListener("htmx:beforeSend", function (event) {
        // console.log(event.target.id)
        // console.log("we received a wsmessage") //troubleshooting 
        qanimation = document.querySelector(".Q-animation-holder"); //
    const target = event.target; // Typecast event.target to HTMLElement
    // if statement to only get target websocket sends from an html item with a card id
    if (target.id && target.id.startsWith("card")) {
        gridItem = target; // Assign the clicked element
        //get the card's grid index from its array. 
        const gridIndex = Array.from(gridContainer.children).indexOf(gridItem);

        const columns = 6; // Number of columns in the grid
        const rows = 5 //Number of rows in the grid

        let columnStart = (Math.ceil((gridIndex+1)/ (rows))); // GridIndex gets a +1 as it starts in 0. Index goes along rows, first column has 0,1,2,3 and 4 index and so on.
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
            qanimation.classList.add("Q-transition")
        }
    }
    });
    document.addEventListener("htmx:beforeRequest",  (event) =>  {
        
        console.log("we received a request")
        qanimation = document.querySelector(".Q-animation-holder")
        const questionCover = document.getElementById("question-cover")
        if (questionCover != null){
            console.log("we received a wsmessage and there is a question cover")
            qanimation.classList.remove("Q-animation")
        }
        
    });
    
    // Listen for the end of the transition
    document.addEventListener("transitionend", function (event) {
        const target = event.target
        if (target.classList.contains("Q-animation-holder") && target.classList.contains("Q-transition") &&
            !target.classList.contains("Q-animation") && event.propertyName === "transform"){
            console.log("end of animation, triggering removal of Q-animation holder")
            htmx.trigger(target, "htmx:animationEnd");
        }
    });   
}