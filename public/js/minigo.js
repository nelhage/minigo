(function () {
   function longColor(c) {
     if(c == "W") {
       return "white";
     } else if(c == "B") {
       return "black";
     } else {
       return '';
     }
   }
   var size = 9;

   var GoSquare = React.createClass({
       doMove: function(e) {
         e.preventDefault();
         this.props.onSubmitMove({x: this.props.x, y: this.props.y})
       },
       render: function() {
           var classes = ["stone"];
           var c = longColor(this.props.contents);
           if (c == "") {
             classes.push("empty");
           } else {
             classes.push(c)
           }
           return (
             <div className="square" data-coords={JSON.stringify([this.props.x,this.props.y])}>
               <div className={classes.join(" ")} onClick={this.doMove}></div>
               <div className="grid"></div>
             </div>
           );
       }
   });

   var GoBoard = React.createClass({
       getInitialState: function() {
         return {
           positions: {},
           to_move: 'W',
         };
       },
       componentDidMount: function() {
         $.ajax({
           url: "/board.json",
           dataType: 'json',
           cache: false,
           success: function(data) {
             this.setState(data);
           }.bind(this),
           error: function(xhr, status, err) {
             console.error("/board.json", status, err.toString());
           }.bind(this),
         });
       },
       at: function(x, y) {
         return this.state.positions[x+","+y];
       },
       submitMove: function(pos){
         $.ajax({
           method: 'POST',
           url: "/move",
           dataType: 'json',
           cache: false,
           data: JSON.stringify({
             x: pos.x,
             y: pos.y,
             to_move: this.state.to_move,
           }),
           success: function(data) {
             this.setState(data);
           }.bind(this),
           error: function(xhr, status, err) {
             console.error("doMove", status, err.toString());
           }.bind(this),
         })
       },
       render: function() {
         var rows = [];
         for (var y = 0; y < this.props.size; y++) {
           var row = [];
           for (var x = 0; x < this.props.size; x++) {
             row.push(
                 <GoSquare
                     key={x} x={x} y={y}
                     contents={this.at(x,y)}
                     onSubmitMove={this.submitMove}
                 />);
           }
           rows.push(<div className="row" key={y}>{row}</div>);
         }
         var classes = ["goboard", longColor(this.state.to_move)];
         return (
           <div className={classes.join(" ")}>
             {rows}
           </div>
         );
       }
   });

   ReactDOM.render(
     <GoBoard size={size} />,
     document.getElementById('content')
   );
 })();
