
const [PAGE_WIDTH, PAGE_HEIGHT] = [595.28, 841.89];

const PAPER_PADDING = 36;

export async function outputPDF({
  element,
  rate,
  pageHeight,
  version
}) {
  if (!(element instanceof HTMLElement)) {
    return;
  }
  const A4_HEIGHT = pageHeight;

  function addBlank(_height, dom, type) {
    // console.log({_height, dom, type})
    let blankDiv = document.createElement('div');
    // blankDiv.style.backgroundColor = 'red';
    blankDiv.style.width = '200px';
    let divHeight = _height;
    if (version <= 102) {
      divHeight = divHeight/rate;
    }
    // addHeight += divHeight;
    blankDiv.style.height = divHeight + 'px';
    dom.dataset.isAdd = "true";
    dom.insertAdjacentElement(type, blankDiv);
    if (dom.classList && dom.classList.contains('el-table__row') && dom.parentElement) {
      let newHeight = dom.parentElement.offsetHeight + divHeight + 'px';
      let findParentElement = findTableParentElement(dom, 'el-table__inner-wrapper');
      if (findParentElement) {
        findParentElement.style.height = newHeight;
        // findParentElement.style.backgroundColor = '#67c23a';
      }
    } else {
      // dom.style.backgroundColor = '#67c23a';
    }
  }

  function findTableParentElement(dom, findKey,count = 0) {
    let findParentElement;
    if (dom && dom.parentElement && (count <= 8)) {
      if (dom.parentElement.classList && dom.parentElement.classList.contains(findKey)) {
        findParentElement = dom.parentElement
      } else {
        findParentElement = findTableParentElement(dom.parentElement, findKey, ++count)
      }
    }
    return findParentElement;
  }

  const originalPageHeight = A4_HEIGHT;


  function getElementTop(contentElement) {
    if (contentElement.getBoundingClientRect) {
      const rect = contentElement.getBoundingClientRect() || {};
      const topDistance = rect.top;
      return topDistance;
    }
  }

  let t = getElementTop(element);

  function traversingNodes(nodes, parentTop = 0) {
    for (const element of nodes) {
      const one = element;

      let isTableRow =
        one.classList && one.classList.contains('el-table__row');
      
      let isTableEmpty =
        one.classList && one.classList.contains('el-table__empty-block');
      
      let isFormItem =
        one.classList && one.classList.contains('el-form-item');

      let isDividerTitle =
        one.classList && one.classList.contains('print-divider-title');
      
      
      let isPageDom = 
        one.classList && one.classList.contains('page-col');

      const { offsetHeight } = one;

      const elementTop = getElementTop(one);


      let top = elementTop - t - parentTop;
      const rateOffsetHeight = offsetHeight;

      let tempMarginTop;
      try {
        tempMarginTop = window.getComputedStyle(one).marginTop;
        tempMarginTop = parseFloat(tempMarginTop);
      } catch (error) {
        tempMarginTop = 0;
      }

      
      if (isPageDom && (version > 102)) {
        updatePagePos(rateOffsetHeight + tempMarginTop*2, one);
      } else if (isTableRow || isFormItem || isTableEmpty) {
        updateTablePos(rateOffsetHeight, top, one);
      } else if (isDividerTitle) {
        updateTablePos(rateOffsetHeight + tempMarginTop*2, top, one);
      } else {
        updateNormalElPos(rateOffsetHeight, top, one, parentTop);
      }
    }
  }


  function updateNormalElPos(eHeight, top, dom, parentTop) {
    if (
      ((top) % originalPageHeight) + eHeight > (originalPageHeight)
    ) {
      if (dom.childNodes) {
        traversingNodes(dom.childNodes, parentTop);
      } else {
        let _height = originalPageHeight - ((top) % originalPageHeight);
        if (!dom.dataset.isAdd) {
          addBlank(_height, dom, 'beforebegin')
        }
      }
    }
  }


  function updateTablePos(eHeight, top, dom) {
    if (
      (((top) % originalPageHeight) + eHeight) > (originalPageHeight)
    ) {
      let _height = originalPageHeight - ((top) % originalPageHeight);
      _height=_height < eHeight?eHeight:_height;
      if (!dom.dataset.isAdd) {
        addBlank(_height, dom, 'beforebegin')
      }
    }
  }

  let pagePosHeight = 0;
  function updatePagePos(eHeight, dom) {
    if (eHeight > (originalPageHeight)) {
      addPageDelimiter(dom)
      dom.style.pageBreakAfter = 'always';
      let top = getElementTop(dom) - t;
      traversingNodes(dom.childNodes, top);
    } else if ((pagePosHeight + eHeight) > (originalPageHeight)) {
      addPageDelimiter(dom)
      pagePosHeight += eHeight;
    } else {
      pagePosHeight += eHeight;
    }
  }

  function addPageDelimiter(dom) {
    dom.style.pageBreakBefore = 'always';
    pagePosHeight = 0;
    // dom.style.backgroundColor = '#67c23a';
  }

  traversingNodes(element.childNodes);

  return;
}
