sea.ui={
    ScreenFill : function(width,height,processing,displayPoint){
	//p
	this.p = processing;
	this.width = width;
	this.height = height;
	this.widgetArray = new Array();
	this.displayPoint = displayPoint;
	
	//m
	this.setWidth = function(w){this.width = w;};
	this.setHeight = function(h){this.height = h;};
	this.addWidget = function(widget){
	    this.widgetArray.push(widget);
	};
	this.addWidgetArray = function(widgetArray){
	    this.widgetArray = this.widgetArray.concat(widgetArray);
	}
	this.draw = function(){
	    for(var i = 0; i < this.widgetArray.length; i++){
		if(i == displayPoint.length){
		    break;
		}
		var point = displayPoint[i];
		this.widgetArray[i].draw(this.p,point.x,point.y);
	    }
	};
	this.hide = function(){
	    alert("screen fill hide");
	};
	this.refresh = function(){
	    alert("screen refresh");
	};
    }
    ,WidgetBase : function(){
	this.draw = function(p,ox,oy,w,h){
	    p.noStroke();
	    p.fill(210,10);
	    p.rectMode(p.CENTER);
	    p.rect(ox,oy,w,h);
	};
    }
    ,WidgetSentence : function(sentence){
	this.sentence = sentence;
	var widget = new sea.ui.WidgetBase();

	this.draw = function(p,ox,oy){
	    // widget.draw(p,ox,oy,150,100);
	    p.fill(50);
	    p.textSize(60);
	    p.textAlign(p.CENTER);
	    p.text(this.sentence.title,ox,oy);
	    //var asccenter = p.textAscent();
	    //var textWidth = p.textWidth(this.sentence.title);
	    //p.line(ox - textWidth/2,oy,ox + textWidth/2,oy);
	    //p.ellipse(ox,oy,10,10);
	};
    }
};
