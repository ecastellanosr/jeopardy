let gridItem: HTMLElement | null = null;
const gridContainer: HTMLElement = document.getElementById('cards-subgrid') ?? document.createElement("div"); // Replace with your target item
const anycard: HTMLElement = document.getElementById('card1') ?? document.createElement("div");
const dimensions: DOMRect = anycard.getBoundingClientRect() ?? document.createElement("div");
const cardWidth = dimensions.width
const cardHeight = dimensions.height


// Calculate left and bottom
const GridWidth = gridContainer.offsetWidth; // Assuming 6 columns
const GridHeight = gridContainer.offsetHeight; // Assuming 6 rows

const WidthRatio = cardWidth/GridWidth
const HeightRatio = cardHeight/GridHeight

// Get the computed styles of the grid item

gridContainer.addEventListener("htmx:afterRequest", function (event) {
  const target = event.target as HTMLElement; // Typecast event.target to HTMLElement
  if (target.id && target.id.startsWith("card")) {
    gridItem = target; // Assign the clicked element
    const qanimation: HTMLElement = document.getElementById("Q-animation")  ?? document.createElement("div");
    const itemStyle = window.getComputedStyle(gridItem);
    // Get the start and end positions
    const columnStart = parseInt(itemStyle.getPropertyValue('grid-column-start'), 10);
    const columnEnd = parseInt(itemStyle.getPropertyValue('grid-column-end'), 10);
    const rowStart = parseInt(itemStyle.getPropertyValue('grid-row-start'), 10);
    const rowEnd = parseInt(itemStyle.getPropertyValue('grid-row-end'), 10);
  
    const columnWidth = gridContainer.offsetWidth/6; // Assuming 6 columns
    const columnHeight = gridContainer.offsetHeight/5; // Assuming 6 rows
  
    const x = (columnStart - 1) * columnWidth;
    const y = gridContainer.offsetHeight - rowEnd * columnHeight;
    qanimation.style.setProperty("--x",`${x}`); 
    qanimation.style.setProperty("--y",`${y}`);
    qanimation.style.setProperty("--width",`${cardWidth} px`);
    qanimation.style.setProperty("--heigth",`${cardHeight} px`);
    qanimation.style.setProperty("--Ratiox",`${WidthRatio}`);
    qanimation.style.setProperty("--Ratioy",`${HeightRatio}`);
  }
});


   


